// the listener connection code has been disabled:
// it would rarely be used anyway since a delay is now added to crawler jobs;
// if re-enabled in the future, some work is needed to gracefully handle disconnection

package server

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type server struct {
	stopped chan struct{}
	query   *dao.Query
	//pool       *pgxpool.Pool
	handlers   []handler.Handler
	gcInterval time.Duration
	logger     *zap.SugaredLogger
}

func (s *server) Start(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()
	//pListenerConn, listenerConnErr := s.newListenerConn(ctx)
	//if listenerConnErr != nil {
	//	err = listenerConnErr
	//	return
	//}
	//listenerConn := pListenerConn.Conn()
	handlers := make([]serverHandler, len(s.handlers))
	listenerChans := make(map[string]chan pgconn.Notification)
	for i, h := range s.handlers {
		listenerChan := make(chan pgconn.Notification)
		sh := serverHandler{
			Handler: h,
			sem:     semaphore.NewWeighted(int64(h.Concurrency)),
			query:   s.query,
			//listenerConn: listenerConn,
			listenerChan: listenerChan,
			logger:       s.logger.With("queue", h.Queue),
		}
		handlers[i] = sh
		listenerChans[h.Queue] = listenerChan
		//if _, listenErr := listenerConn.Exec(ctx, fmt.Sprintf(`LISTEN %q`, h.Queue)); listenErr != nil {
		//	err = listenErr
		//	return
		//}
		go sh.start(ctx)
	}
	go func() {
		for {
			select {
			case <-s.stopped:
				cancel()
			case <-ctx.Done():
				//pListenerConn.Release()
				return
			}
		}
	}()
	//go func() {
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			return
	//		default:
	//			notification, waitErr := listenerConn.WaitForNotification(ctx)
	//			if waitErr != nil {
	//				if !errors.Is(waitErr, context.Canceled) {
	//					s.logger.Errorf("Error waiting for notification: %s", waitErr)
	//				}
	//				continue
	//			}
	//			ch, ok := listenerChans[notification.Channel]
	//			if !ok {
	//				s.logger.Errorf("Received notification for unknown channel: %s", notification.Channel)
	//				continue
	//			}
	//			select {
	//			case <-ctx.Done():
	//				return
	//			case ch <- *notification:
	//				continue
	//			}
	//		}
	//	}
	//}()
	go s.runGarbageCollection(ctx)
	return
}

//func (s *server) newListenerConn(ctx context.Context) (*pgxpool.Conn, error) {
//	conn, err := s.pool.Acquire(ctx)
//	if err != nil {
//		return nil, err
//	}
//	_, err = conn.Exec(ctx, "SET idle_in_transaction_session_timeout = 0")
//	if err != nil {
//		return nil, err
//	}
//	return conn, nil
//}

func (s *server) runGarbageCollection(ctx context.Context) {
	for {
		tx := s.query.QueueJob.WithContext(ctx).Where(
			s.query.QueueJob.Status.In(string(model.QueueJobStatusProcessed), string(model.QueueJobStatusFailed)),
		).UnderlyingDB().Where("queue_jobs.ran_at + queue_jobs.archival_duration < ?::timestamptz", time.Now()).Delete(&model.QueueJob{})
		if tx.Error != nil {
			s.logger.Errorw("error deleting old queue jobs", "error", tx.Error)
		} else if tx.RowsAffected > 0 {
			s.logger.Debugw("deleted old queue jobs", "count", tx.RowsAffected)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(s.gcInterval):
			continue
		}
	}
}

type serverHandler struct {
	handler.Handler
	sem   *semaphore.Weighted
	query *dao.Query
	//listenerConn *pgx.Conn
	listenerChan chan pgconn.Notification
	logger       *zap.SugaredLogger
}

func (h *serverHandler) start(ctx context.Context) {
	checkTicker := time.NewTicker(1)
	for {
		select {
		case <-ctx.Done():
			return
		case notification := <-h.listenerChan:
			if semErr := h.sem.Acquire(ctx, 1); semErr != nil {
				return
			}
			go func() {
				defer h.sem.Release(1)
				_, _, _ = h.handleJob(ctx, h.query.QueueJob.ID.Eq(notification.Payload))
			}()
		case <-checkTicker.C:
			if semErr := h.sem.Acquire(ctx, 1); semErr != nil {
				return
			}
			checkTicker.Reset(h.CheckInterval)
			go func() {
				defer h.sem.Release(1)
				jobId, _, err := h.handleJob(ctx)
				// if a job was found, we should check straight away for another job, otherwise we wait for the check interval
				if err == nil && jobId != "" {
					checkTicker.Reset(1)
				}
			}()
		}
	}
}

func (h *serverHandler) handleJob(ctx context.Context, conds ...gen.Condition) (jobId string, processed bool, err error) {
	err = h.query.Transaction(func(tx *dao.Query) error {
		job, findErr := tx.QueueJob.WithContext(ctx).Where(
			append(conds,
				h.query.QueueJob.Queue.Eq(h.Queue),
				h.query.QueueJob.Status.In(string(model.QueueJobStatusPending), string(model.QueueJobStatusRetry)),
				h.query.QueueJob.RunAfter.Lte(time.Now()),
			)...,
		).Order(
			h.query.QueueJob.Status.Eq(string(model.QueueJobStatusRetry)),
			h.query.QueueJob.Priority,
			h.query.QueueJob.RunAfter,
		).Clauses(clause.Locking{
			Strength: "UPDATE",
			Options:  "SKIP LOCKED",
		}).First()
		if findErr != nil {
			if errors.Is(findErr, gorm.ErrRecordNotFound) {
				return nil
			}
			return findErr
		}
		jobId = job.ID
		var jobErr error
		if job.Deadline.Valid && job.Deadline.Time.Before(time.Now()) {
			jobErr = ErrJobExceededDeadline
			h.logger.Debugw("job deadline is in the past, skipping", "job_id", job.ID)
		} else {
			// check if the job is being retried and increment retry count accordingly
			if job.Status != model.QueueJobStatusPending {
				job.Retries++
			}
			// execute the queue handler of this job
			jobErr = handler.Exec(ctx, h.Handler, *job)
		}

		job.RanAt = sql.NullTime{Time: time.Now(), Valid: true}

		if jobErr != nil {
			h.logger.Errorw("job failed", "error", jobErr)
			if job.Retries < job.MaxRetries {
				job.Status = model.QueueJobStatusRetry
				job.RunAfter = queue.CalculateBackoff(job.Retries)
			} else {
				job.Status = model.QueueJobStatusFailed
			}
			job.Error = model.NewNullString(jobErr.Error())
		} else {
			job.Status = model.QueueJobStatusProcessed
			processed = true
		}
		_, updateErr := tx.QueueJob.WithContext(ctx).Updates(job)
		return updateErr
	})
	if err != nil {
		h.logger.Error("error handling job", "error", err)
	} else if processed {
		h.logger.Debugw("job processed", "job_id", jobId)
	}
	return
}

var ErrJobExceededDeadline = errors.New("the job did not complete before its deadline")
