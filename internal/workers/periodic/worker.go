package periodic

import (
	"context"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

type worker struct {
	fn              func(context.Context) error
	interval        time.Duration
	initialInterval time.Duration
	quickShutdown   bool
}

func New(interval time.Duration, fn func(context.Context) error, options ...Option) runner.Provider {
	wrk := &worker{
		fn:              fn,
		interval:        interval,
		initialInterval: interval,
	}

	for _, opt := range options {
		opt(wrk)
	}

	return wrk
}

func (w *worker) Runner() runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		shutdown := make(chan struct{})

		// initialInterval := w.initialInterval
		// if initialInterval == 0 {
		// 	initialInterval = time.Millisecond
		// }
		// interval := w.initialInterval

		// ticker := time.NewTicker(w.initialInterval)

		// initial := true
		wait := time.After(w.initialInterval)

		go func() {
			defer cancel(nil)

			for {
				select {
				case <-shutdown:
					return
				case <-wait:
					err := w.fn(ctx)
					if err != nil {
						cancel(fmt.Errorf("%w: %w: %w", Err, ErrInvoke, err))
						return
					}

					wait = time.After(w.interval)
				}
			}
		}()

		return func(shutdownCtx context.Context) error {
			close(shutdown)

			if w.quickShutdown {
				cancel(runner.ErrShutdownRequested)
			} else {
				select {
				case <-shutdownCtx.Done():
					return fmt.Errorf("%w: %w: %w", Err, ErrShutdown, shutdownCtx.Err())
				case <-ctx.Done():
				}
			}

			return nil
		}, nil
	}
}
