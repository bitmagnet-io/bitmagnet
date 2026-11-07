package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const (
	DefaultHandlerTimeout = 30 * time.Second
)

// Func is a function that Handlers execute for every Job on a queue
type Func func(ctx context.Context, job model.QueueJob) error

// Handler handles jobs on a queue
type Handler struct {
	Queue         string
	Concurrency   int
	CheckInterval time.Duration
	JobTimeout    time.Duration
	Handle        Func
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
		Handle:        f,
		Queue:         queue,
		CheckInterval: 30 * time.Second,
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

// Exec runs a queue job with a real deadline. Cancellation must reach the
// handler (and therefore MatchLab/TMDB requests), not just the caller waiting
// for the result. A buffered result channel cannot strand a finished goroutine
// after a timeout.
func Exec(ctx context.Context, handler Handler, job model.QueueJob) error {
	deadlineCtx, cancel := context.WithTimeout(ctx, handler.JobTimeout)
	defer cancel()

	result := make(chan error, 1)
	go func() {
		var runErr error
		defer func() {
			if panicValue := recover(); panicValue != nil {
				log.Printf("recovering from a panic in the job handler:\n%s", debug.Stack())
				_, file, line, ok := runtime.Caller(1)
				if ok && strings.Contains(file, "runtime/") {
					_, file, line, ok = runtime.Caller(2)
				}
				if ok {
					runErr = fmt.Errorf("panic [%s:%d]: %v", file, line, panicValue)
				} else {
					runErr = fmt.Errorf("panic: %v", panicValue)
				}
			}
			result <- runErr
		}()
		runErr = handler.Handle(deadlineCtx, job)
	}()

	select {
	case err := <-result:
		if err != nil {
			return fmt.Errorf("job failed to process: %w", err)
		}
		return nil
	case <-deadlineCtx.Done():
		if errors.Is(deadlineCtx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("job exceeded its %s timeout: %w", handler.JobTimeout, deadlineCtx.Err())
		}
		return deadlineCtx.Err()
	}
}
