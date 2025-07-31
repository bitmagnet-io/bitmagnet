package channel

import (
	"context"
	"fmt"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

type Adder[T any] interface {
	Add(context.Context, ...T) error
}

type Worker[T any] interface {
	runner.Provider
	Adder[T]
}

func NewWorker[T any](fn func(context.Context, T) error, options ...Option[T]) Worker[T] {
	wrk := &worker[T]{
		size: 1,
		fn:   fn,
	}

	for _, opt := range options {
		opt(wrk)
	}

	return wrk
}

type worker[T any] struct {
	mtx           sync.RWMutex
	size          int
	fn            func(context.Context, T) error
	ch            chan []T
	shutdown      chan struct{}
	quickShutdown bool
}

func (w *worker[T]) Add(ctx context.Context, items ...T) error {
	if len(items) == 0 {
		return nil
	}

	w.mtx.RLock()
	ch, shutdown := w.ch, w.shutdown
	w.mtx.RUnlock()

	if ch == nil {
		return fmt.Errorf("%w: %w", Err, ErrUninitialized)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-shutdown:
		return fmt.Errorf("%w: %w", Err, runner.ErrShutdownRequested)
	case ch <- items:
		return nil
	}
}

func (w *worker[T]) Runner() runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		w.mtx.Lock()
		defer w.mtx.Unlock()

		if w.ch != nil {
			return runner.NopShutdowner, runner.ErrAlreadyRunning
		}

		ch := make(chan []T)
		shutdown := make(chan struct{})
		sem := make(chan struct{}, w.size)

		w.ch = ch
		w.shutdown = shutdown

		go func() {
			defer func() {
				for range w.size {
					sem <- struct{}{}
				}

				cancel(nil)
			}()

			for {
				select {
				case <-shutdown:
					return
				case items, ok := <-ch:
					if !ok {
						return
					}

					for _, item := range items {
						select {
						case <-ctx.Done():
							return
						case sem <- struct{}{}:
							go func(item T) {
								defer func() {
									<-sem
								}()

								err := w.fn(ctx, item)
								if err != nil {
									cancel(fmt.Errorf("%w: %w: %w", Err, ErrItem, err))
								}
							}(item)
						}
					}
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

			w.mtx.Lock()
			w.ch = nil
			w.shutdown = nil
			w.mtx.Unlock()

			close(ch)

			return nil
		}, nil
	}
}
