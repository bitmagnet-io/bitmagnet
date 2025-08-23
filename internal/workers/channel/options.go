package channel

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/metrics"
)

type Option[T any] func(*worker[T])

func WithSize[T any](size int) Option[T] {
	return WithAtomicSize[T](atomic.NewValue(size))
}

func WithAtomicSize[T any](size *atomic.Value[int]) Option[T] {
	return func(wrk *worker[T]) {
		wrk.size = size
	}
}

func WithQuickShutdown[T any]() Option[T] {
	return func(wrk *worker[T]) {
		wrk.quickShutdown = true
	}
}

func WithMetricsAdapter[T any](adapter metrics.Adapter) Option[T] {
	return func(wrk *worker[T]) {
		wrk.metrics = adapter
	}
}

func WithOnIdle[T any](onIdle func(context.Context) error) Option[T] {
	return func(wrk *worker[T]) {
		wrk.onIdle = onIdle
	}
}
