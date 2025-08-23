package channel

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/semaphore"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

type Adder[T any] interface {
	Add(context.Context, ...T) error
}

type Worker[T any] interface {
	runner.Provider
	Adder[T]
}

type Func[T any] func(context.Context, T) error

func NewWorker[T any](fn Func[T], options ...Option[T]) Worker[T] {
	wrk := &worker[T]{
		fn:      fn,
		metrics: metrics.Nop,
		onIdle: func(context.Context) error {
			return nil
		},
	}

	for _, opt := range options {
		opt(wrk)
	}

	if wrk.size == nil {
		wrk.size = atomic.NewValue(1)
	}

	return wrk
}

type worker[T any] struct {
	mtx           sync.RWMutex
	size          *atomic.Value[int]
	fn            Func[T]
	ch            chan []T
	shutdown      chan struct{}
	quickShutdown bool
	metrics       metrics.Adapter
	onIdle        func(context.Context) error
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
		sem, semUnsubscribe := semaphore.NewAtomic(w.size)

		w.ch = ch
		w.shutdown = shutdown

		go func() {
			defer func() {
				sem.Acquire(ctx, semUnsubscribe())

				cancel(nil)
				w.metrics.Reset()
			}()

			for {
				select {
				case <-shutdown:
					return
				case items, ok := <-ch:
					if !ok {
						return
					}

					timeAdded := time.Now()

					w.metrics.IncrAdded(len(items))

					for _, item := range items {
						if err := sem.Acquire(ctx, 1); err != nil {
							return
						}
						go func(item T) {
							defer func() {
								if n := sem.Release(1); n == 1 {
									select {
									case <-shutdown:
										return
									default:
										if err := w.onIdle(ctx); err != nil {
											cancel(err)
										}
									}
								}
							}()

							w.metrics.IncrDequeued(time.Since(timeAdded))

							err := w.fn(ctx, item)
							if err != nil {
								cancel(fmt.Errorf("%w: %w: %w", Err, ErrItem, err))

								return
							}

							w.metrics.IncrFlushed(time.Since(timeAdded))
						}(item)
					}
				}
			}
		}()

		return func(shutdownCtx context.Context) error {
			w.mtx.Lock()
			defer w.mtx.Unlock()

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

			w.ch = nil
			w.shutdown = nil

			close(ch)

			return nil
		}, nil
	}
}
