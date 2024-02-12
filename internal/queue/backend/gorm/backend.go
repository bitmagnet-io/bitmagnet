package gorm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// DefaultConnectionTimeout defines the default amount of time that Queue waits for connections to become available.
	DefaultConnectionTimeout = 30 * time.Second
)

// backend is a Postgres-based Queue backend
type backend struct {
	config         queue.Config               // backend configuration
	cancelFuncs    []context.CancelFunc       // cancel functions to be called upon Shutdown()
	handlers       map[string]handler.Handler // a map of queue names to queue handlers
	newQueues      chan string                // a channel that indicates that new queues are ready to be processed
	readyQueues    chan string                // a channel that indicates which queues are ready to have jobs processed.
	listenCancelCh chan context.CancelFunc    // cancellation channel for the listenerConn's WaitForNotification call.
	listenerConn   *pgx.Conn                  // dedicated connection that LISTENs for jobs across all queues
	listenerConnMu sync.Mutex                 // listenerConnMu protects the listener connection from concurrent access
	mu             sync.Mutex                 // protects concurrent access to fields on backend
	query          *dao.Query
	pool           *pgxpool.Pool // connection pool for backend, used to process and enqueue jobs
	logger         *zap.SugaredLogger
}

// newQueueMonitor monitors for new queues and instruct's the listener connection to LISTEN for jobs on them
func (b *backend) newQueueMonitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case newQueue := <-b.newQueues:
			b.logger.Debug("configure new handler", "queue", newQueue)
		setupListeners:
			// drain b.listenCancelCh before setting up new listeners
			select {
			case cancelListener := <-b.listenCancelCh:
				b.logger.Debug("canceling previous wait listeners", "queue", newQueue)
				cancelListener()
				goto setupListeners
			default:
			}

			b.listenerConnMu.Lock()
			if b.listenerConn == nil {
				if lc, err := b.newListenerConn(ctx); err != nil {
					b.logger.Errorw("unable to initialize listener connection", "queue", newQueue, "error", err)
					b.listenerConnMu.Unlock()
					return
				} else {
					b.listenerConn = lc
				}
			}
			// note: 'LISTEN, channel' is idempotent
			_, err := b.listenerConn.Exec(ctx, fmt.Sprintf(`LISTEN %q`, newQueue))
			b.listenerConnMu.Unlock()
			if err != nil {
				err = fmt.Errorf("unable to configure listener connection: %w", err)
				b.logger.Errorw("FATAL ERROR unable to listen for new jobs", "queue", newQueue, "error", err)
				return
			}

			b.logger.Debug("listening on queue", "queue", newQueue)
			b.readyQueues <- newQueue
		}
	}
}

func (b *backend) newListenerConn(ctx context.Context) (*pgx.Conn, error) {
	conn, err := b.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	_, err = conn.Exec(ctx, "SET idle_in_transaction_session_timeout = 0")
	if err != nil {
		return nil, err
	}
	return conn.Conn(), nil
}

// Enqueue adds jobs to the specified queue
func (b *backend) Enqueue(ctx context.Context, job model.QueueJob) (jobID string, err error) {
	if err := b.query.QueueJob.WithContext(ctx).Create(&job); err != nil {
		return "", err
	}
	return job.ID, nil
}

// Start starts processing jobs with the specified queue and handler
func (b *backend) Start(ctx context.Context, h handler.Handler) (err error) {
	ctx, cancel := context.WithCancel(ctx)

	b.logger.Debugw("starting job processing", "queue", h.Queue)
	b.mu.Lock()
	b.cancelFuncs = append(b.cancelFuncs, cancel)
	b.handlers[h.Queue] = h
	b.mu.Unlock()

	b.newQueues <- h.Queue

	err = b.start(ctx, h)
	if err != nil {
		b.logger.Errorw("unable to start processing queue", "queue", h.Queue, "error", err)
		return
	}
	return
}

// Shutdown shuts this backend down
func (b *backend) Shutdown(context.Context) {
	b.logger.Debug("starting shutdown")

	for _, f := range b.cancelFuncs {
		f()
	}

	b.cancelFuncs = nil
	b.logger.Debug("shutdown complete")
}

// start starts processing new, pending, and future jobs
func (b *backend) start(ctx context.Context, h handler.Handler) (err error) {
	var ok bool
	var listenJobChan chan *pgconn.Notification
	var errCh chan error

	if h, ok = b.handlers[h.Queue]; !ok {
		return fmt.Errorf("%w: %s", handler.ErrNoHandlerForQueue, h.Queue)
	}

	pendingJobsChan := b.pendingJobs(ctx, h.Queue)

	// wait for the listener to connect and be ready to listen
	for q := range b.readyQueues {
		if q == h.Queue {
			listenJobChan, errCh = b.listen(ctx)
			break
		}

		b.logger.Debugw("Picked up a queue that a different start() will be waiting for. Adding back to ready list",
			"queue", q)
		b.readyQueues <- q
	}

	for i := 0; i < h.Concurrency; i++ {
		go func() {
			var err error
			var n *pgconn.Notification

			for {
				select {
				case n = <-listenJobChan:
					err = b.handleJob(ctx, n.Payload)
				case n = <-pendingJobsChan:
					err = b.handleJob(ctx, n.Payload)
				case <-ctx.Done():
					return
				case <-errCh:
					b.logger.Error("error hanlding job", "error", err)
					continue
				}

				if err != nil {
					if errors.Is(err, context.Canceled) {
						err = nil
						continue
					}

					b.logger.Errorw(
						"job failed",
						"queue", h.Queue,
						"error", err,
						"job_id", n.Payload,
					)

					continue
				}
			}
		}()
	}

	return nil
}

func (b *backend) pendingJobs(ctx context.Context, queue string) (jobsCh chan *pgconn.Notification) {
	jobsCh = make(chan *pgconn.Notification)

	go func(ctx context.Context) {

		for {
			jobID, err := b.getPendingJobID(ctx, queue)
			if err == nil && jobID != "" {
				select {
				case <-ctx.Done():
					return
				case jobsCh <- &pgconn.Notification{Channel: queue, Payload: jobID}:
					continue
				}
			} else {
				if err != nil && !errors.Is(err, context.Canceled) {
					b.logger.Errorw(
						"failed to fetch pending job",
						"queue", queue,
						"error", err,
						"job_id", jobID,
					)
				}
				select {
				case <-ctx.Done():
					return
				case <-time.After(b.config.JobCheckInterval):
					continue
				}
			}
		}
	}(ctx)

	return jobsCh
}

// handleJob is the workhorse of Queue
// it receives pending, periodic, and retry job ids asynchronously
// 1. handleJob first creates a transactions inside of which a row lock is acquired for the job to be processed.
// 2. handleJob secondly calls the handler on the job, and finally updates the job's status
func (b *backend) handleJob(ctx context.Context, jobID string) error {
	return b.query.Transaction(func(tx *dao.Query) error {
		job, findErr := b.query.QueueJob.WithContext(ctx).Where(
			b.query.QueueJob.ID.Eq(jobID),
			b.query.QueueJob.Status.In(string(model.QueueJobStatusPending), string(model.QueueJobStatusRetry)),
			b.query.QueueJob.RunAfter.Lte(time.Now()),
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
		var jobErr error
		if job.Deadline.Valid && job.Deadline.Time.Before(time.Now()) {
			jobErr = queue.ErrJobExceededDeadline
			b.logger.Debugw("job deadline is in the past, skipping", "queue", job.Queue, "job_id", job.ID)
		} else {
			// check if the job is being retried and increment retry count accordingly
			if job.Status != queue.JobStatusPending {
				job.Retries++
			}
			h, ok := b.handlers[job.Queue]
			if !ok {
				b.logger.Errorw("received a job for which no handler is configured",
					"queue", job.Queue,
					"job_id", job.ID)
				return handler.ErrNoHandlerForQueue
			}
			// execute the queue handler of this job
			jobErr = handler.Exec(ctx, h, *job)
		}

		job.RanAt = sql.NullTime{Time: time.Now(), Valid: true}

		if jobErr != nil {
			b.logger.Errorw("job failed", "error", jobErr)
			if job.Retries < job.MaxRetries {
				job.Status = queue.JobStatusRetry
				job.RunAfter = queue.CalculateBackoff(job.Retries)
			} else {
				job.Status = queue.JobStatusFailed
			}
			job.Error = model.NewNullString(jobErr.Error())
		} else {
			job.Status = queue.JobStatusProcessed
		}
		_, updateErr := b.query.QueueJob.WithContext(ctx).Updates(job)
		return updateErr
	})
}

// listen uses Postgres LISTEN to listen for jobs on a queue
// TODO: There is currently no handling of listener disconnects in backend.
// This will lead to jobs not getting processed until the worker is restarted.
// Implement disconnect handling.
func (b *backend) listen(ctx context.Context) (c chan *pgconn.Notification, errCh chan error) {
	c = make(chan *pgconn.Notification)
	errCh = make(chan error)

	waitForNotificationCtx, cancel := context.WithCancel(ctx)
	b.listenCancelCh <- cancel

	go func(ctx context.Context) {
		var notification *pgconn.Notification
		var waitErr error
		for {
			select {
			case <-ctx.Done():
				// our context has been canceled, the system is shutting down
				return
			default:
				b.listenerConnMu.Lock()
				notification, waitErr = b.listenerConn.WaitForNotification(waitForNotificationCtx)
				b.listenerConnMu.Unlock()
			}
			if waitErr != nil {
				if errors.Is(waitErr, context.Canceled) {
					// this is likely not a system shutdown, but an interrupt from the goroutine that manages changes to
					// the list of handlers. It needs the connection to be unbusy so that it can instruct the connection
					// to start listening on any new queues
					b.logger.Debug("Stopping notifications processing")
					return
				}

				// The connection is busy adding new LISTENers
				if b.listenerConn.PgConn().IsBusy() {
					b.logger.Debug("listen connection is busy, trying to acquire listener connection again...")
					waitForNotificationCtx, cancel = context.WithCancel(ctx)
					b.listenCancelCh <- cancel
					continue
				}

				b.logger.Errorw("failed to wait for notification", "error", waitErr)
				continue
			}

			b.logger.Debugw(
				"job notification for queue",
				"notification", notification,
				"err", waitErr,
			)

			c <- notification
		}
	}(ctx)

	return c, errCh
}

func (b *backend) getPendingJobID(ctx context.Context, queue string) (jobID string, err error) {
	err = b.query.QueueJob.WithContext(ctx).Where(
		b.query.QueueJob.Queue.Eq(queue),
		b.query.QueueJob.Status.In(string(model.QueueJobStatusPending), string(model.QueueJobStatusRetry)),
		b.query.QueueJob.RunAfter.Lte(time.Now()),
	).Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "SKIP LOCKED",
	}).Limit(1).Pluck(b.query.QueueJob.ID, &jobID)
	return
}
