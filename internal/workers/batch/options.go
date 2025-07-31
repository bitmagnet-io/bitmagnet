package batch

import (
	"sync"
	"time"
)

type Option[K comparable, V any] func(*worker[K, V])

func WithKeyer[K comparable, V any](keyer Keyer[K, V]) Option[K, V] {
	return func(wrk *worker[K, V]) {
		wrk.keyer = keyer
	}
}

func WithFilter[K comparable, V any](filter Filter[K, V]) Option[K, V] {
	return func(wrk *worker[K, V]) {
		currentFilter := wrk.filter

		wrk.filter = func(k K, v V) bool {
			return currentFilter(k, v) && filter(k, v)
		}
	}
}

func WithMerger[K comparable, V any](merger Merger[V]) Option[K, V] {
	return func(wrk *worker[K, V]) {
		wrk.merger = merger
	}
}

func WithFlusher[K comparable, V any](flusher Flusher[V]) Option[K, V] {
	return func(wrk *worker[K, V]) {
		wrk.flusher = flusher
	}
}

func WithMaxSize[K comparable, V any](maxSize int) Option[K, V] {
	return func(wrk *worker[K, V]) {
		wrk.maxSize = maxSize
	}
}

func WithMaxWait[K comparable, V any](maxWait time.Duration) Option[K, V] {
	return func(wrk *worker[K, V]) {
		wrk.maxWait = maxWait
	}
}

func WithValuesAsKeys[K comparable]() Option[K, K] {
	return WithKeyer(func(key K) K {
		return key
	})
}

func WithSequentialKeys[V any]() Option[int, V] {
	var (
		mtx sync.Mutex
		key int
	)

	return WithKeyer(func(V) int {
		mtx.Lock()
		defer mtx.Unlock()

		key++

		return key
	})
}

func WithQuickShutdown[K comparable, V any]() Option[K, V] {
	return func(wrk *worker[K, V]) {
		wrk.quickShutdown = true
	}
}
