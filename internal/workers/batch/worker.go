package batch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
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
	runner.Provider
	Adder[V]
}

func NewWorker[K comparable, V any](options ...Option[K, V]) Worker[V] {
	wrk := &worker[K, V]{
		keyer:   defaultKeyer[K, V],
		filter:  defaultFilter[K, V],
		merger:  defaultMerger[V],
		flusher: defaultFlusher[V],
		maxSize: 1,
		metrics: metrics.Nop,
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
	metrics       metrics.Adapter
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
		return fmt.Errorf("%w: %w", Err, runner.ErrShutdownRequested)
	case ch <- items:
		return nil
	}
}

func (w *worker[K, V]) Runner() runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		w.mtx.Lock()
		defer w.mtx.Unlock()

		if w.ch != nil {
			return runner.NopShutdowner, runner.ErrAlreadyRunning
		}

		w.ch = make(chan []V, 1)
		w.shutdown = make(chan struct{})

		wr := workerRunner[K, V]{
			worker:    w,
			ctx:       ctx,
			cancel:    cancel,
			buffer:    maps.NewInsertMap[K, value[V]](),
			lastFlush: time.Now(),
		}

		go wr.run()

		return wr.shutdowner, nil
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
