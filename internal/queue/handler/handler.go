package handler

import (
	"runtime"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

const (
	DefaultHandlerTimeout = 30 * time.Second
)

// Func is a function that Handlers execute for every Job on a queue
type Func func(job model.QueueJob) runner.Runner

// Handler handles jobs on a queue
type Handler struct {
	Queue         string
	Concurrency   int
	CheckInterval time.Duration
	JobTimeout    time.Duration
	Func          Func
}

// Option is function that sets optional configuration for Handlers
type Option func(w *Handler)

// WithOptions sets one or more options on handler
func (h *Handler) WithOptions(opts ...Option) {
	for _, opt := range opts {
		opt(h)
	}
}

// JobTimeout configures handlers with a time deadline for every executed job
// The timeout is the amount of time that can be spent executing the handler's Func
// when a timeout is exceeded, the job fails and enters its retry phase
func JobTimeout(d time.Duration) Option {
	return func(h *Handler) {
		h.JobTimeout = d
	}
}

// Concurrency configures Queue handlers to process jobs concurrently
// the default concurrency is the number of (v)CPUs on the machine running Queue
func Concurrency(c int) Option {
	return func(h *Handler) {
		h.Concurrency = c
	}
}

// New creates new queue handlers for specific queues. This function is to be usued to create new Handlers for
// non-periodic jobs (most jobs). Use [NewPeriodic] to initialize handlers for periodic jobs.
func New(queue string, f Func, opts ...Option) (h Handler) {
	h = Handler{
		Func:          f,
		Queue:         queue,
		CheckInterval: 10 * time.Second,
	}

	h.WithOptions(opts...)

	// default to running on as many goroutines as there are CPUs
	if h.Concurrency == 0 {
		Concurrency(runtime.NumCPU())(&h)
	}

	// always set a job timeout if none is set
	if h.JobTimeout == 0 {
		JobTimeout(DefaultHandlerTimeout)(&h)
	}

	return
}

// Exec executes handler functions with a timeout
// func Exec(ctx context.Context, handler Handler, job model.QueueJob) (worker.Shutdowner, error) {
//	var (
//		shutdowner worker.Shutdowner
//		errs       []error
//		mtx        sync.Mutex
//	)
//
//	addErr := func(e error) {
//		mtx.Lock()
//		defer mtx.Unlock()
//		errs = append(errs, e)
//	}
//
//	//errCh := make(chan error, 1)
//	done := make(chan struct{})
//
//	go func(ctx context.Context) {
//		defer func() {
//			if x := recover(); x != nil {
//				log.Printf("recovering from a panic in the job handler:\n%s", string(debug.Stack()))
//
//				_, file, line, ok := runtime.Caller(1) // skip the first frame (panic itself)
//				if ok && strings.Contains(file, "runtime/") {
//					// The panic came from the runtime, most likely due to incorrect
//					// map/slice usage. The parent frame should have the real trigger.
//					_, file, line, ok = runtime.Caller(2)
//				}
//
//				// Include the file and line number info in the error, if runtime.Caller returned ok.
//				if ok {
//					addErr(fmt.Errorf("panic [%s:%d]: %v", file, line, x))
//				} else {
//					addErr(fmt.Errorf("panic: %v", x))
//				}
//			}
//
//			close(done)
//		}()
//
//		runner := handler.Func(job)
//
//		var err error
//
//		ctx, cancel := context.WithTimeout(ctx, handler.JobTimeout)
//		defer cancel()
//
//		runCtx, runCancel := context.CancelC
//
//		mtx.Lock()
//		shutdowner, err = runner(runCtx, runCancel)
//	}(ctx)
//
//	select {
//	case <-done:
//		err = <-errCh
//		if err != nil {
//			err = fmt.Errorf("job failed to process: %w", err)
//		}
//
//	case <-timeoutCtx.Done():
//		ctxErr := timeoutCtx.ErrConcat()
//
//		switch {
//		case errors.Is(ctxErr, context.DeadlineExceeded):
//			err = fmt.Errorf("job exceeded its %s timeout: %w", handler.JobTimeout, ctxErr)
//		case errors.Is(ctxErr, context.Canceled):
//			err = ctxErr
//		default:
//			err = fmt.Errorf("job failed to process: %w", ctxErr)
//		}
//	}
//
//	return err
//}
