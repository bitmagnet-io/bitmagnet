package server

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"go.uber.org/zap"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type server struct {
	daoProvider database.DaoTransactionProvider
	handlers    []handler.Handler
	gcInterval  time.Duration
	gcSemaphore chan struct{}
	draining    chan struct{}
	logger      *zap.SugaredLogger
}

func (s *server) Start(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
	handlers := make([]serverHandler, len(s.handlers))

	for i, h := range s.handlers {
		sh := serverHandler{
			Handler:     h,
			sem:         make(chan struct{}, h.Concurrency),
			draining:    s.draining,
			daoProvider: s.daoProvider,
			logger:      s.logger.With("queue", h.Queue),
		}

		handlers[i] = sh

		go sh.start(ctx)
	}

	go s.runGarbageCollection(ctx)

	return func(ctx context.Context) error {
		defer cancel(runner.ErrShutdownRequested)

		close(s.draining)

		var wg sync.WaitGroup

		wg.Add(len(handlers))

		for _, h := range handlers {
			go func(h serverHandler) {
				defer wg.Done()
				h.drain(ctx)
			}(h)
		}

		s.gcSemaphore <- struct{}{}

		wg.Wait()

		return nil
	}, nil
}

func (s *server) runGarbageCollection(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.draining:
			return
		case s.gcSemaphore <- struct{}{}:
		}

		daoQ, err := s.daoProvider.Dao()
		if err != nil {
			s.logger.Errorw("error getting dao", "error", err)
		}

		tx := daoQ.QueueJob.WithContext(ctx).Where(
			daoQ.QueueJob.Status.In(string(model.QueueJobStatusProcessed), string(model.QueueJobStatusFailed)),
		).
			UnderlyingDB().Where(
			"queue_jobs.ran_at + queue_jobs.archival_duration < ?::timestamptz",
			time.Now(),
		).Delete(&model.QueueJob{})

		if tx.Error != nil {
			s.logger.Errorw("error deleting old queue jobs", "error", tx.Error)
		} else if tx.RowsAffected > 0 {
			s.logger.Debugw("deleted old queue jobs", "count", tx.RowsAffected)
		}

		<-s.gcSemaphore

		select {
		case <-ctx.Done():
			return
		case <-s.draining:
			return
		case <-time.After(s.gcInterval):
			continue
		}
	}
}

type serverHandler struct {
	handler.Handler
	sem         chan struct{}
	draining    chan struct{}
	daoProvider database.DaoTransactionProvider
	// listenerChan chan pgconn.Notification
	logger *zap.SugaredLogger
}

func (h *serverHandler) start(ctx context.Context) {
	checkTicker := time.NewTicker(1)

	for {
		select {
		case <-ctx.Done():
			return
		case <-h.draining:
			return
		case <-checkTicker.C:
			select {
			case <-ctx.Done():
				return
			case <-h.draining:
				return
			case h.sem <- struct{}{}:
			}

			checkTicker.Reset(h.CheckInterval)

			go func() {
				jobID, _, err := h.handleJob(ctx)
				// if a job was found, we should check straight away for another job,
				// otherwise we wait for the check interval
				if err == nil && jobID != "" {
					checkTicker.Reset(1)
				}

				<-h.sem
			}()
		}
	}
}

func (h *serverHandler) handleJob(
	ctx context.Context,
	conds ...gen.Condition,
) (jobID string, processed bool, err error) {
	err = h.daoProvider.DaoTransaction(func(tx *dao.Query) error {
		job, findErr := tx.QueueJob.WithContext(ctx).Where(
			append(
				conds,
				tx.QueueJob.Queue.Eq(h.Queue),
				tx.QueueJob.Status.In(
					string(model.QueueJobStatusPending),
					string(model.QueueJobStatusRetry),
				),
				tx.QueueJob.RunAfter.Lte(time.Now()),
			)...,
		).Order(
			tx.QueueJob.Status.Eq(string(model.QueueJobStatusRetry)),
			tx.QueueJob.Priority,
			tx.QueueJob.RunAfter,
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

		jobID = job.ID

		var jobErr error
		if job.Deadline.Valid && job.Deadline.Time.Before(time.Now()) {
			jobErr = ErrJobExceededDeadline

			h.logger.Debugw("job deadline is in the past, skipping", "job_id", job.ID)
		} else {
			// check if the job is being retried and increment retry count accordingly
			if job.Status != model.QueueJobStatusPending {
				job.Retries++
			}

			jobErr = h.executeJob(ctx, *job)
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
		h.logger.Errorw("error handling job", "error", err, "job_id", jobID)
	} else if processed {
		h.logger.Debugw("job processed", "job_id", jobID)
	}

	return
}

func (h *serverHandler) executeJob(ctx context.Context, job model.QueueJob) error {
	run := h.Func(job)

	runCtx, runCancel := context.WithCancelCause(ctx)

	shutdowner, err := run(runCtx, runCancel)
	if err != nil {
		runCancel(err)
		return errors.Join(err, shutdowner(ctx))
	}

	select {
	case <-runCtx.Done():
		return extractContextCause(runCtx)
	case <-time.After(h.JobTimeout):
		runCancel(ErrJobExceededTimeout)
		return ErrJobExceededTimeout
	case <-h.draining:
		return shutdowner(runCtx)
	}
}

func extractContextCause(ctx context.Context) error {
	cause := context.Cause(ctx)

	if errors.Is(cause, runner.ErrCompleted) {
		return nil
	}

	return cause
}

func (h *serverHandler) drain(ctx context.Context) {
	for range h.Concurrency {
		select {
		case <-ctx.Done():
			return
		case h.sem <- struct{}{}:
		}
	}
}

var (
	ErrJobExceededDeadline = errors.New("the job did not complete before its deadline")
	ErrJobExceededTimeout  = errors.New("the job exceeded its timeout")
)
