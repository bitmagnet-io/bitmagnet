package batch

import (
	"context"
	"fmt"
	"sync"
	"time"

	runner2 "github.com/bitmagnet-io/bitmagnet/internal/workers/runner"

	"github.com/bitmagnet-io/bitmagnet/internal/maps"
)

type Keyer[K comparable, V any] func(V) K

type HasKey[K comparable] interface {
	Key() K
}

type Filter[K comparable, V any] func(K, V) bool

type Merger[V any] func(V, V) V

type Flusher[V any] func(context.Context, []V) error

type Adder[V any] interface {
	Add(context.Context, ...V) error
}

type Worker[V any] interface {
	runner2.Provider
	Adder[V]
}

func NewWorker[K comparable, V any](options ...Option[K, V]) Worker[V] {
	wrk := &worker[K, V]{
		keyer:   defaultKeyer[K, V],
		filter:  defaultFilter[K, V],
		merger:  defaultMerger[V],
		flusher: defaultFlusher[V],
		maxSize: 1,
	}

	for _, opt := range options {
		opt(wrk)
	}

	return wrk
}

type worker[K comparable, V any] struct {
	mtx           sync.RWMutex
	ch            chan []V
	shutdown      chan struct{}
	keyer         Keyer[K, V]
	merger        Merger[V]
	flusher       Flusher[V]
	filter        Filter[K, V]
	maxSize       int
	maxWait       time.Duration
	quickShutdown bool
}

func (w *worker[K, V]) Add(ctx context.Context, items ...V) error {
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
		return fmt.Errorf("%w: %w", Err, runner2.ErrShutdownRequested)
	case ch <- items:
		return nil
	}
}

func (w *worker[K, V]) Runner() runner2.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner2.Shutdowner, error) {
		w.mtx.Lock()
		defer w.mtx.Unlock()

		if w.ch != nil {
			return runner2.NopShutdowner, runner2.ErrAlreadyRunning
		}

		ch := make(chan []V, 1)
		shutdown := make(chan struct{})
		buff := maps.NewInsertMap[K, V]()
		lastFlush := time.Now()

		w.ch = ch
		w.shutdown = shutdown

		addItemsAndCheckFlush := func(ctx context.Context, items []V, forceFlush bool) error {
			w.mtx.Lock()

			for _, item := range items {
				key := w.keyer(item)

				if !w.filter(key, item) {
					continue
				}

				if existing, ok := buff.Get(key); ok {
					item = w.merger(existing, item)
				}
				buff.Set(key, item)
			}

			size := buff.Len()

			shouldFlush := size > 0 && (forceFlush || size > w.maxSize || time.Since(lastFlush) > w.maxWait)

			if shouldFlush {
				items = buff.Values()
				buff.Clear()
			}

			lastFlush = time.Now()

			w.mtx.Unlock()

			if shouldFlush {
				err := w.flusher(ctx, items)
				if err != nil {
					return fmt.Errorf("%w: %w: %w", Err, ErrFlush, err)
				}
			}

			return nil
		}

		go func() {
			defer cancel(nil)

			for {
				select {
				case <-shutdown:
					_ = addItemsAndCheckFlush(ctx, nil, true)
					return
				case <-time.After(w.maxWait):
					_ = addItemsAndCheckFlush(ctx, nil, true)
				case items, ok := <-ch:
					if !ok {
						return
					}

					err := addItemsAndCheckFlush(ctx, items, false)
					if err != nil {
						cancel(err)
						return
					}
				}
			}
		}()

		return func(shutdownCtx context.Context) error {
			close(shutdown)

			var err error

			if w.quickShutdown {
				cancel(runner2.ErrShutdownRequested)
			} else {
				select {
				case <-shutdownCtx.Done():
					err = fmt.Errorf("%w: %w: %w", Err, ErrShutdown, shutdownCtx.Err())
				case <-ctx.Done():
				}
			}

			w.mtx.Lock()
			w.ch = nil
			w.shutdown = nil
			w.mtx.Unlock()

			close(ch)

			return err
		}, nil
	}
}

func defaultKeyer[K comparable, V any](v V) K {
	k, ok := any(v).(K)

	if ok {
		return k
	}

	if hasKey, ok := any(v).(HasKey[K]); ok {
		return hasKey.Key()
	}

	return k
}

func defaultMerger[V any](a, b V) V {
	return b
}

func defaultFlusher[V any](context.Context, []V) error {
	return nil
}

func defaultFilter[K comparable, V any](K, V) bool {
	return true
}
