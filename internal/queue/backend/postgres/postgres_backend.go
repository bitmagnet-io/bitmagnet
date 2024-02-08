package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/acaloiaro/neoq"
	"github.com/acaloiaro/neoq/handler"
	"github.com/acaloiaro/neoq/jobs"
	"github.com/acaloiaro/neoq/logging"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jsuar/go-cron-descriptor/pkg/crondescriptor"
	"github.com/robfig/cron"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

const (
	JobQuery = `SELECT id,fingerprint,queue,status,deadline,payload,retries,max_retries,run_after,ran_at,created_at,error
					FROM neoq_jobs
					WHERE id = $1
					AND status NOT IN ('processed')
					FOR UPDATE SKIP LOCKED
					LIMIT 1`
	PendingJobIDQuery = `SELECT id
					FROM neoq_jobs
					WHERE queue = $1
					AND status NOT IN ('processed')
					AND run_after <= NOW()
					FOR UPDATE SKIP LOCKED
					LIMIT 1`
	FutureJobQuery = `SELECT id,fingerprint,queue,status,deadline,payload,retries,max_retries,run_after,ran_at,created_at,error
					FROM neoq_jobs
					WHERE queue = $1
					AND status NOT IN ('processed')
					AND run_after > NOW()
					ORDER BY run_after ASC
					LIMIT 100
					FOR UPDATE SKIP LOCKED`
	setIdleInTxSessionTimeout = `SET idle_in_transaction_session_timeout = 0`
	pgConnectionBusyRetries   = 10 // the number of times to retry busy postgres connections, i.e. PgConn.IsBusy()
)

type txContextKey struct{}

var (
	// DefaultConnectionTimeout defines the default amount of time that Neoq waits for connections to become available.
	DefaultConnectionTimeout         = 30 * time.Second
	txCtxVarKey                      txContextKey
	shutdownJobID                    = "-1" // job ID announced when triggering a shutdown
	shutdownAnnouncementAllowance    = 100  // ms
	ErrCnxString                     = errors.New("invalid connecton string: see documentation for valid connection strings")
	ErrDuplicateJob                  = errors.New("duplicate job")
	ErrNoTransactionInContext        = errors.New("context does not have a Tx set")
	ErrExceededConnectionPoolTimeout = errors.New("exceeded timeout acquiring a connection from the pool")
)

// PgBackend is a Postgres-based Neoq backend
type PgBackend struct {
	neoq.Neoq
	cancelFuncs    []context.CancelFunc       // cancel functions to be called upon Shutdown()
	config         *neoq.Config               // backend configuration
	cron           *cron.Cron                 // scheduler for periodic jobs
	futureJobs     map[string]*jobs.Job       // map of future job IDs to the corresponding job record
	handlers       map[string]handler.Handler // a map of queue names to queue handlers
	newQueues      chan string                // a channel that indicates that new queues are ready to be processed
	readyQueues    chan string                // a channel that indicates which queues are ready to have jobs processed.
	listenCancelCh chan context.CancelFunc    // cancellation channel for the listenerConn's WaitForNotification call.
	listenerConn   *pgx.Conn                  // dedicated connection that LISTENs for jobs across all queues
	listenerConnMu *sync.RWMutex              // listenerConnMu protects the listener connection from concurrent access
	logger         logging.Logger             // backend-wide logger
	mu             *sync.RWMutex              // protects concurrent access to fields on PgBackend
	pool           *pgxpool.Pool              // connection pool for backend, used to process and enqueue jobs
}

// Backend initializes a new postgres-backed neoq backend
//
// If the database does not yet exist, Neoq will attempt to create the database and related tables by default.
//
// Backend requires that one of the [neoq.ConfigOption] is [WithConnectionString]
//
// Connection strings may be a URL or DSN-style connection strings. The connection string supports multiple
// options detailed below.
//
// options:
//   - pool_max_conns: integer greater than 0
//   - pool_min_conns: integer 0 or greater
//   - pool_max_conn_lifetime: duration string
//   - pool_max_conn_idle_time: duration string
//   - pool_health_check_period: duration string
//   - pool_max_conn_lifetime_jitter: duration string
//
// # Example DSN
//
// user=worker password=secret host=workerdb.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10
//
// # Example URL
//
// postgres://worker:secret@workerdb.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
func Backend(ctx context.Context, opts ...neoq.ConfigOption) (pb neoq.Neoq, err error) {
	cfg := neoq.NewConfig()
	cfg.IdleTransactionTimeout = neoq.DefaultIdleTxTimeout
	cfg.PGConnectionTimeout = DefaultConnectionTimeout

	p := &PgBackend{
		cancelFuncs:    []context.CancelFunc{},
		config:         cfg,
		cron:           cron.New(),
		futureJobs:     make(map[string]*jobs.Job),
		handlers:       make(map[string]handler.Handler),
		newQueues:      make(chan string),
		readyQueues:    make(chan string),
		listenerConnMu: &sync.RWMutex{},
		mu:             &sync.RWMutex{},
		listenCancelCh: make(chan context.CancelFunc, 1),
	}

	// Set all options
	for _, opt := range opts {
		opt(p.config)
	}

	p.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: p.config.LogLevel}))
	ctx, cancel := context.WithCancel(ctx)
	p.mu.Lock()
	p.cancelFuncs = append(p.cancelFuncs, cancel)
	p.mu.Unlock()

	if p.pool == nil { //nolint: nestif
		var poolConfig *pgxpool.Config
		poolConfig, err = pgxpool.ParseConfig(p.config.ConnectionString)
		if err != nil || p.config.ConnectionString == "" {
			return nil, ErrCnxString
		}

		// ensure that workers don't consume connections with idle transactions
		poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) (err error) {
			var query string
			if p.config.IdleTransactionTimeout > 0 {
				query = fmt.Sprintf("SET idle_in_transaction_session_timeout = '%dms'", p.config.IdleTransactionTimeout)
			} else {
				// there is no limit to the amount of time a worker's transactions may be idle
				query = setIdleInTxSessionTimeout
			}

			if !p.config.SynchronousCommit {
				query = fmt.Sprintf("%s; SET synchronous_commit = 'off';", query)
			}
			_, err = conn.Exec(ctx, query)
			return
		}

		p.pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			return
		}
	}

	p.listenerConn, err = p.newListenerConn(ctx)
	if err != nil {
		p.logger.Error("unable to initialize listener connection", slog.Any("error", err))
	}

	// monitor handlers for changes and LISTEN when new queues are added
	go p.newQueueMonitor(ctx)

	p.cron.Start()

	pb = p

	return pb, nil
}

// newQueueMonitor monitors for new queues and instruct's the listener connection to LISTEN for jobs on them
func (p *PgBackend) newQueueMonitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case newQueue := <-p.newQueues:
			p.logger.Debug("configure new handler", "queue", newQueue)
		setup_listeners:
			// drain p.listenCancelCh before setting up new listeners
			select {
			case cancelListener := <-p.listenCancelCh:
				p.logger.Debug("canceling previous wait listeners", "queue", newQueue)
				cancelListener()
				goto setup_listeners
			default:
			}

			p.listenerConnMu.Lock()
			// note: 'LISTEN, channel' is idempotent
			_, err := p.listenerConn.Exec(ctx, fmt.Sprintf(`LISTEN %q`, newQueue))
			p.listenerConnMu.Unlock()
			if err != nil {
				err = fmt.Errorf("unable to configure listener connection: %w", err)
				p.logger.Error("FATAL ERROR unable to listen for new jobs", slog.String("queue", newQueue), slog.Any("error", err))
				return
			}

			p.logger.Debug("listening on queue", "queue", newQueue)
			p.readyQueues <- newQueue
		}
	}
}

func (p *PgBackend) newListenerConn(ctx context.Context) (conn *pgx.Conn, err error) {
	var pgxCfg *pgx.ConnConfig
	pgxCfg, err = pgx.ParseConfig(p.config.ConnectionString)
	if err != nil {
		return
	}

	// remove any pgxpool parameters before creating a new connection
	customPgxParams := []string{
		"pool_max_conns", "pool_min_conns",
		"pool_max_conn_lifetime", "pool_max_conn_idle_time", "pool_health_check_period",
		"pool_max_conn_lifetime_jitter",
	}
	for param := range pgxCfg.RuntimeParams {
		if slices.Contains(customPgxParams, param) {
			delete(pgxCfg.RuntimeParams, param)
		}
	}
	conn, err = pgx.ConnectConfig(ctx, pgxCfg)
	if err != nil {
		p.logger.Error("unable to acquire listener connection", slog.Any("error", err))
		return
	}
	_, err = conn.Exec(ctx, "SET idle_in_transaction_session_timeout = 0")

	return
}

// WithConnectionString configures neoq postgres backend to use the specified connection string when connecting to a backend
func WithConnectionString(connectionString string) neoq.ConfigOption {
	return func(c *neoq.Config) {
		c.ConnectionString = connectionString
	}
}

// WithTransactionTimeout sets the time that PgBackend's transactions may be idle before its underlying connection is
// closed
// The timeout is the number of milliseconds that a transaction may sit idle before postgres terminates the
// transaction's underlying connection. The timeout should be longer than your longest job takes to complete. If set
// too short, job state will become unpredictable, e.g. retry counts may become incorrect.
func WithTransactionTimeout(txTimeout int) neoq.ConfigOption {
	return func(c *neoq.Config) {
		c.IdleTransactionTimeout = txTimeout
	}
}

// WithConnectionTimeout sets the duration that Neoq waits for connections to become available to process and enqueue jobs
//
// Note: ConnectionTimeout does not affect how long neoq waits for connections to run schema migrations
func WithConnectionTimeout(timeout time.Duration) neoq.ConfigOption {
	return func(c *neoq.Config) {
		c.PGConnectionTimeout = timeout
	}
}

// WithSynchronousCommit enables postgres parameter `synchronous_commit`.
//
// By default, neoq runs with synchronous_commit disabled.
//
// Postgres incurrs significant transactional overhead from synchronously committing small transactions. Because
// neoq jobs must be enqueued individually, and payloads are generally quite small, synchronous_commit introduces
// significant overhead, but increases data durability.
//
// See https://www.postgresql.org/docs/current/wal-async-commit.html for details on the implications that this has for
// neoq jobs.
//
// Enabling synchronous commit results in an order of magnitude slowdown in enqueueing and processing jobs.
func WithSynchronousCommit(enabled bool) neoq.ConfigOption {
	return func(c *neoq.Config) {
		c.SynchronousCommit = enabled
	}
}

// txFromContext gets the transaction from a context, if the transaction is already set
func txFromContext(ctx context.Context) (t pgx.Tx, err error) {
	var ok bool
	if t, ok = ctx.Value(txCtxVarKey).(pgx.Tx); ok {
		return
	}

	err = ErrNoTransactionInContext

	return
}

// Enqueue adds jobs to the specified queue
func (p *PgBackend) Enqueue(ctx context.Context, job *jobs.Job) (jobID string, err error) {
	if job.Queue == "" {
		err = jobs.ErrNoQueueSpecified
		return
	}

	p.logger.Debug("enqueueing job payload", slog.String("queue", job.Queue), slog.Any("job_payload", job.Payload))

	p.logger.Debug("acquiring new connection from connection pool", slog.String("queue", job.Queue))
	conn, err := p.acquire(ctx)
	if err != nil {
		err = fmt.Errorf("error acquiring connection: %w", err)
		return
	}
	defer conn.Release()

	p.logger.Debug("beginning new transaction to enqueue job", slog.String("queue", job.Queue))
	tx, err := conn.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("error creating transaction: %w", err)
		return
	}

	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer func(ctx context.Context) { _ = tx.Rollback(ctx) }(ctx) // rollback has no effect if the transaction has been committed
	jobID, err = p.enqueueJob(ctx, tx, job)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				err = ErrDuplicateJob
				return
			}
		}
		p.logger.Error("error enqueueing job", slog.String("queue", job.Queue), slog.Any("error", err))
		err = fmt.Errorf("error enqueuing job: %w", err)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		err = fmt.Errorf("error committing transaction: %w", err)
		return
	}
	p.logger.Debug("job added to queue:", slog.String("queue", job.Queue), slog.String("job_id", jobID))

	// add future jobs to the future job list
	if job.RunAfter.After(time.Now().UTC()) {
		p.mu.Lock()
		p.futureJobs[jobID] = job
		p.mu.Unlock()
		p.logger.Debug(
			"added job to future jobs list",
			slog.String("queue", job.Queue),
			slog.String("job_id", jobID),
			slog.Time("run_after", job.RunAfter),
		)
	}

	return jobID, nil
}

// Start starts processing jobs with the specified queue and handler
func (p *PgBackend) Start(ctx context.Context, h handler.Handler) (err error) {
	ctx, cancel := context.WithCancel(ctx)

	p.logger.Debug("starting job processing", slog.String("queue", h.Queue))
	p.mu.Lock()
	p.cancelFuncs = append(p.cancelFuncs, cancel)
	p.handlers[h.Queue] = h
	p.mu.Unlock()

	p.newQueues <- h.Queue

	err = p.start(ctx, h)
	if err != nil {
		p.logger.Error("unable to start processing queue", slog.String("queue", h.Queue), slog.Any("error", err))
		return
	}
	return
}

// StartCron starts processing jobs with the specified cron schedule and handler
//
// See: https://pkg.go.dev/github.com/robfig/cron?#hdr-CRON_Expression_Format for details on the cron spec format
func (p *PgBackend) StartCron(ctx context.Context, cronSpec string, h handler.Handler) (err error) {
	cd, err := crondescriptor.NewCronDescriptor(cronSpec)
	if err != nil {
		p.logger.Error(
			"error creating cron descriptor",
			slog.String("queue", h.Queue),
			slog.String("cronspec", cronSpec),
			slog.Any("error", err),
		)
		return fmt.Errorf("error creating cron descriptor: %w", err)
	}

	cdStr, err := cd.GetDescription(crondescriptor.Full)
	if err != nil {
		p.logger.Error(
			"error getting cron descriptor",
			slog.String("queue", h.Queue),
			slog.Any("descriptor", crondescriptor.Full),
			slog.Any("error", err),
		)
		return fmt.Errorf("error getting cron description: %w", err)
	}

	queue := StripNonAlphanum(strcase.ToSnake(*cdStr))
	h.Queue = queue

	ctx, cancel := context.WithCancel(ctx)
	p.mu.Lock()
	p.cancelFuncs = append(p.cancelFuncs, cancel)
	p.mu.Unlock()

	if err = p.cron.AddFunc(cronSpec, func() {
		_, err := p.Enqueue(ctx, &jobs.Job{Queue: queue})
		if err != nil {
			// When we are working with a cron we want to ignore the canceled and the duplicate job errors. The duplicate job
			// error specifically is not one the cron enqueuer needs to concern itself with because that means that another
			// worker has already enqueued the job for this cron recurrence. It is not helpful to log the error in that
			// scenario since the job will be processed.
			if errors.Is(err, context.Canceled) || errors.Is(err, ErrDuplicateJob) {
				return
			}

			p.logger.Error("error queueing cron job", slog.String("queue", h.Queue), slog.Any("error", err))
		}
	}); err != nil {
		return fmt.Errorf("error adding cron: %w", err)
	}

	return p.Start(ctx, h)
}

// SetLogger sets this backend's logger
func (p *PgBackend) SetLogger(logger logging.Logger) {
	p.logger = logger
}

// Shutdown shuts this backend down
func (p *PgBackend) Shutdown(ctx context.Context) {
	p.logger.Debug("starting shutdown")
	for queue := range p.handlers {
		p.announceJob(ctx, queue, shutdownJobID)
	}

	// wait for the announcement to process
	time.Sleep(time.Duration(shutdownAnnouncementAllowance) * time.Millisecond)

	for _, f := range p.cancelFuncs {
		f()
	}

	p.pool.Close()
	p.cron.Stop()

	p.cancelFuncs = nil
	p.logger.Debug("shutdown complete")
}

// enqueueJob adds jobs to the queue, returning the job ID
//
// Jobs that are not already fingerprinted are fingerprinted before being added
// Duplicate jobs are not added to the queue. Any two unprocessed jobs with the same fingerprint are duplicates
func (p *PgBackend) enqueueJob(ctx context.Context, tx pgx.Tx, j *jobs.Job) (jobID string, err error) {
	err = jobs.FingerprintJob(j)
	if err != nil {
		return
	}

	p.logger.Debug("adding job to the queue", slog.String("queue", j.Queue))
	err = tx.QueryRow(ctx, `INSERT INTO neoq_jobs(queue, fingerprint, payload, run_after, deadline, max_retries)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		j.Queue, j.Fingerprint, j.Payload, j.RunAfter, j.Deadline, j.MaxRetries).Scan(&jobID)
	if err != nil {
		err = fmt.Errorf("unable add job to queue: %w", err)
		return
	}

	return jobID, err
}

// moveToDeadQueue moves jobs from the pending queue to the dead queue
func (p *PgBackend) moveToDeadQueue(ctx context.Context, tx pgx.Tx, j *jobs.Job, jobErr string) (err error) {
	_, err = tx.Exec(ctx, "DELETE FROM neoq_jobs WHERE id = $1", j.ID)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, `INSERT INTO neoq_dead_jobs(id, queue, fingerprint, payload, retries, max_retries, error, deadline)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		j.ID, j.Queue, j.Fingerprint, j.Payload, j.Retries, j.MaxRetries, jobErr, j.Deadline)

	return
}

// updateJob updates the status of jobs with: status, run time, error messages, and retries
//
// if the retry count exceeds the maximum number of retries for the job, move the job to the dead jobs queue
//
// if `tx`'s underlying connection dies while updating job status, the transaction will fail, and the job's original
// status will be reflecting in the database.
//
// The implication of this is that:
// - the job's 'error' field will not reflect any errors the occurred in the handler
// - the job's retry count is not incremented
// - the job's run time will remain its original value
// - the job has its original 'status'
//
// ultimately, this means that any time a database connection is lost while updating job status, then the job will be
// processed at least one more time.
// nolint: cyclop
func (p *PgBackend) updateJob(ctx context.Context, jobErr error) (err error) {
	status := JobStatusProcessed
	errMsg := ""

	if jobErr != nil {
		p.logger.Error("job failed", slog.Any("job_error", jobErr))
		status = JobStatusFailed
		errMsg = jobErr.Error()
	}

	var job *jobs.Job
	if job, err = jobs.FromContext(ctx); err != nil {
		return fmt.Errorf("error getting job from context: %w", err)
	}

	var tx pgx.Tx
	if tx, err = txFromContext(ctx); err != nil {
		return fmt.Errorf("error getting tx from context: %w", err)
	}

	if job.MaxRetries != nil && job.Retries >= *job.MaxRetries {
		err = p.moveToDeadQueue(ctx, tx, job, errMsg)
		return
	}

	var runAfter time.Time
	if status == JobStatusFailed {
		runAfter = CalculateBackoff(job.Retries)
		qstr := "UPDATE neoq_jobs SET ran_at = $1, error = $2, status = $3, retries = $4, run_after = $5 WHERE id = $6"
		_, err = tx.Exec(ctx, qstr, time.Now().UTC(), errMsg, status, job.Retries, runAfter, job.ID)
	} else {
		qstr := "UPDATE neoq_jobs SET ran_at = $1, error = $2, status = $3 WHERE id = $4"
		_, err = tx.Exec(ctx, qstr, time.Now().UTC(), errMsg, status, job.ID)
	}

	if err != nil {
		return
	}

	if time.Until(runAfter) > 0 {
		p.mu.Lock()
		p.futureJobs[fmt.Sprint(job.ID)] = job
		p.mu.Unlock()
	}

	return nil
}

// start starts processing new, pending, and future jobs
// nolint: cyclop
func (p *PgBackend) start(ctx context.Context, h handler.Handler) (err error) {
	var ok bool
	var listenJobChan chan *pgconn.Notification
	var errCh chan error

	if h, ok = p.handlers[h.Queue]; !ok {
		return fmt.Errorf("%w: %s", handler.ErrNoHandlerForQueue, h.Queue)
	}

	pendingJobsChan := p.pendingJobs(ctx, h.Queue) // process overdue jobs *at startup*

	// wait for the listener to connect and be ready to listen
	for q := range p.readyQueues {
		if q == h.Queue {
			listenJobChan, errCh = p.listen(ctx)
			break
		}

		p.logger.Debug("Picked up a queue that a different start() will be waiting for. Adding back to ready list",
			slog.String("queue", q))
		p.readyQueues <- q
	}

	// process all future jobs and retries
	go func() { p.scheduleFutureJobs(ctx, h.Queue) }()

	for i := 0; i < h.Concurrency; i++ {
		go func() {
			var err error
			var n *pgconn.Notification

			for {
				select {
				case n = <-listenJobChan:
					err = p.handleJob(ctx, n.Payload)
				case n = <-pendingJobsChan:
					err = p.handleJob(ctx, n.Payload)
				case <-ctx.Done():
					return
				case <-errCh:
					p.logger.Error("error hanlding job", "error", err)
					continue
				}

				if err != nil {
					if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, context.Canceled) {
						err = nil
						continue
					}

					p.logger.Error(
						"job failed",
						slog.String("queue", h.Queue),
						slog.Any("error", err),
						slog.String("job_id", n.Payload),
					)

					continue
				}
			}
		}()
	}

	return nil
}

// initFutureJobs is intended to be run once to initialize the list of future jobs that must be monitored for
// execution. it should be run only during system startup.
func (p *PgBackend) initFutureJobs(ctx context.Context, queue string) (err error) {
	rows, err := p.pool.Query(ctx, FutureJobQuery, queue)
	if err != nil {
		p.logger.Error("failed to fetch future jobs list", slog.String("queue", queue), slog.Any("error", err))
		return
	}

	futureJobs, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[jobs.Job])
	if err != nil {
		return
	}

	for _, job := range futureJobs {
		p.mu.Lock()
		p.futureJobs[fmt.Sprintf("%d", job.ID)] = job
		p.mu.Unlock()
	}

	return
}

// scheduleFutureJobs announces future jobs using NOTIFY on an interval
func (p *PgBackend) scheduleFutureJobs(ctx context.Context, queue string) {
	err := p.initFutureJobs(ctx, queue)
	if err != nil {
		return
	}

	// check for new future jobs on an interval
	ticker := time.NewTicker(p.config.JobCheckInterval)

	for {
		// loop over list of future jobs, scheduling goroutines to wait for jobs that are due within the next 30 seconds
		p.mu.Lock()
		for jobID, job := range p.futureJobs {
			timeUntillRunAfter := time.Until(job.RunAfter)
			if timeUntillRunAfter <= p.config.FutureJobWindow {
				delete(p.futureJobs, jobID)
				go func(jid string, j *jobs.Job) {
					scheduleCh := time.After(timeUntillRunAfter)
					<-scheduleCh
					p.announceJob(ctx, j.Queue, jid)
				}(jobID, job)
			}
		}
		p.mu.Unlock()

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

// announceJob announces jobs to queue listeners.
//
// When jobs are inserted into the neoq_jobs table, a trigger announces the new job's arrival. This function is to be
// used for announcing jobs that have not been recently inserted into the neoq_jobs table.
//
// Announced jobs are executed by the first worker to respond to the announcement.
func (p *PgBackend) announceJob(ctx context.Context, queue, jobID string) {
	conn, err := p.acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return
	}

	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer func(ctx context.Context) { _ = tx.Rollback(ctx) }(ctx)

	// notify listeners that a job is ready to run
	_, err = tx.Exec(ctx, fmt.Sprintf(`SELECT pg_notify('%s', '%s')`, queue, jobID))
	if err != nil {
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		return
	}
}

func (p *PgBackend) pendingJobs(ctx context.Context, queue string) (jobsCh chan *pgconn.Notification) {
	jobsCh = make(chan *pgconn.Notification)

	conn, err := p.acquire(ctx)
	if err != nil {
		p.logger.Error(
			"failed to acquire database connection to listen for pending queue items",
			slog.String("queue", queue),
			slog.Any("error", err),
		)
		return
	}

	go func(ctx context.Context) {
		defer conn.Release()

		for {
			jobID, err := p.getPendingJobID(ctx, conn, queue)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, context.Canceled) {
					break
				}

				p.logger.Error(
					"failed to fetch pending job",
					slog.String("queue", queue),
					slog.Any("error", err),
					slog.String("job_id", jobID),
				)
			} else {
				jobsCh <- &pgconn.Notification{Channel: queue, Payload: jobID}
			}
		}
	}(ctx)

	return jobsCh
}

// handleJob is the workhorse of Neoq
// it receives pending, periodic, and retry job ids asynchronously
// 1. handleJob first creates a transactions inside of which a row lock is acquired for the job to be processed.
// 2. handleJob secondly calls the handler on the job, and finally updates the job's status
// nolint: cyclop
func (p *PgBackend) handleJob(ctx context.Context, jobID string) (err error) {
	var job *jobs.Job
	var tx pgx.Tx
	conn, err := p.acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	tx, err = conn.Begin(ctx)
	if err != nil {
		return
	}
	defer func(ctx context.Context) { _ = tx.Rollback(ctx) }(ctx) // rollback has no effect if the transaction has been committed

	job, err = p.getJob(ctx, tx, jobID)
	if err != nil {
		return
	}

	if job.Deadline != nil && job.Deadline.Before(time.Now().UTC()) {
		err = jobs.ErrJobExceededDeadline
		p.logger.Debug("job deadline is in the past, skipping", slog.String("queue", job.Queue), slog.Int64("job_id", job.ID))
		err = p.updateJob(ctx, err)
		return
	}

	ctx = withJobContext(ctx, job)
	ctx = context.WithValue(ctx, txCtxVarKey, tx)

	// check if the job is being retried and increment retry count accordingly
	if job.Status != JobStatusNew {
		job.Retries++
	}

	var jobErr error
	h, ok := p.handlers[job.Queue]
	if !ok {
		p.logger.Error("received a job for which no handler is configured",
			slog.String("queue", job.Queue),
			slog.Int64("job_id", job.ID))
		return handler.ErrNoHandlerForQueue
	}

	// execute the queue handler of this job
	jobErr = handler.Exec(ctx, h)
	err = p.updateJob(ctx, jobErr)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		err = fmt.Errorf("error updating job status: %w", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		errMsg := "unable to commit job transaction. retrying this job may dupliate work:"
		p.logger.Error(errMsg, slog.String("queue", h.Queue), slog.Any("error", err), slog.Int64("job_id", job.ID))
		return fmt.Errorf("%s %w", errMsg, err)
	}

	return nil
}

// listen uses Postgres LISTEN to listen for jobs on a queue
// TODO: There is currently no handling of listener disconnects in PgBackend.
// This will lead to jobs not getting processed until the worker is restarted.
// Implement disconnect handling.
func (p *PgBackend) listen(ctx context.Context) (c chan *pgconn.Notification, errCh chan error) {
	c = make(chan *pgconn.Notification)
	errCh = make(chan error)

	waitForNotificationCtx, cancel := context.WithCancel(ctx)
	p.listenCancelCh <- cancel

	go func(ctx context.Context) {
		var notification *pgconn.Notification
		var waitErr error
		for {
			select {
			case <-ctx.Done():
				// our context has been canceled, the system is shutting down
				return
			default:
				p.listenerConnMu.Lock()
				notification, waitErr = p.listenerConn.WaitForNotification(waitForNotificationCtx)
				p.listenerConnMu.Unlock()
			}
			if waitErr != nil {
				if errors.Is(waitErr, context.Canceled) {
					// this is likely not a system shutdown, but an interrupt from the goroutine that manages changes to
					// the list of handlers. It needs the connection to be unbusy so that it can instruct the connection
					// to start listening on any new queues
					p.logger.Debug("Stopping notifications processing")
					return
				}

				// The connection is busy adding new LISTENers
				if p.listenerConn.PgConn().IsBusy() {
					p.logger.Debug("listen connection is busy, trying to acquire listener connection again...")
					waitForNotificationCtx, cancel = context.WithCancel(ctx)
					p.listenCancelCh <- cancel
					continue
				}

				p.logger.Error("failed to wait for notification", slog.Any("error", waitErr))
				continue
			}

			p.logger.Debug(
				"job notification for queue",
				slog.Any("notification", notification),
				slog.Any("err", waitErr),
			)

			// check if Shutdown() has been called
			if notification.Payload == shutdownJobID {
				return
			}

			c <- notification
		}
	}(ctx)

	return c, errCh
}

func (p *PgBackend) getJob(ctx context.Context, tx pgx.Tx, jobID string) (job *jobs.Job, err error) {
	row, err := tx.Query(ctx, JobQuery, jobID)
	if err != nil {
		return
	}

	job, err = pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[jobs.Job])
	if err != nil {
		return
	}

	return
}

func (p *PgBackend) getPendingJobID(ctx context.Context, conn *pgxpool.Conn, queue string) (jobID string, err error) {
	err = conn.QueryRow(ctx, PendingJobIDQuery, queue).Scan(&jobID)
	return
}

// acquire acquires connections from the connection pool with a timeout
//
// the purpose of this function is to skirt pgxpool's default blocking behavior with connection acquisition preemption
func (p *PgBackend) acquire(ctx context.Context) (conn *pgxpool.Conn, err error) {
	ctx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(p.config.PGConnectionTimeout))
	defer cancelFunc()

	p.logger.Debug("acquiring connection with timeout", slog.Any("timeout", p.config.PGConnectionTimeout))

	connCh := make(chan *pgxpool.Conn)
	errCh := make(chan error)

	go func() {
		c, err := p.pool.Acquire(ctx)
		if err != nil {
			errCh <- err
		}

		connCh <- c
	}()

	select {
	case conn = <-connCh:
		return conn, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		p.logger.Error("exceeded timeout acquiring a connection from the pool", slog.Any("timeout", p.config.PGConnectionTimeout))
		cancelFunc()
		err = ErrExceededConnectionPoolTimeout
		return
	}
}

// withJobContext creates a new context with the Job set
func withJobContext(ctx context.Context, j *jobs.Job) context.Context {
	return context.WithValue(ctx, JobCtxVarKey, j)
}
