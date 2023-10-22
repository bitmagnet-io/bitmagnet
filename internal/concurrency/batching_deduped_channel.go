package concurrency

import (
	"context"
	"sync"
	"time"
)

type BatchingDedupedChannel[T any, C comparable] interface {
	In() chan<- T
	InContext(ctx context.Context, value T) bool
	TryIn(value T) bool
	Out() <-chan []T
}

type batchingDedupedChannel[T any, C comparable] struct {
	mutex        sync.Mutex
	toComparable func(T) C
	input        chan T
	output       chan []T
	buffer       map[C]T
	maxBatchSize int
	maxWaitTime  time.Duration
}

func NewBatchingDedupedChannel[T any, C comparable](
	input chan T,
	maxBatchSize int,
	maxWaitTime time.Duration,
	toComparable func(T) C,
) BatchingDedupedChannel[T, C] {
	ch := &batchingDedupedChannel[T, C]{
		toComparable: toComparable,
		input:        input,
		output:       make(chan []T),
		buffer:       make(map[C]T, maxBatchSize),
		maxBatchSize: maxBatchSize,
		maxWaitTime:  maxWaitTime,
	}
	go ch.batch()
	return ch
}

func (ch *batchingDedupedChannel[T, C]) In() chan<- T {
	return ch.input
}

func (ch *batchingDedupedChannel[T, C]) InContext(ctx context.Context, value T) bool {
	select {
	case <-ctx.Done():
		return false
	case ch.input <- value:
		return true
	}
}

func (ch *batchingDedupedChannel[T, C]) TryIn(value T) bool {
	select {
	case ch.input <- value:
		return true
	default:
		return false
	}
}

func (ch *batchingDedupedChannel[T, C]) Out() <-chan []T {
	return ch.output
}

func (ch *batchingDedupedChannel[T, C]) batch() {
	var next T
	var ok bool

	defer close(ch.output)

	for {
		select {
		case next, ok = <-ch.input:
			if !ok {
				break
			}
			ch.mutex.Lock()
			ch.buffer[ch.toComparable(next)] = next
			if len(ch.buffer) >= ch.maxBatchSize {
				ch.flush()
			}
			ch.mutex.Unlock()
		case <-time.After(ch.maxWaitTime):
			ch.mutex.Lock()
			if len(ch.buffer) >= 0 {
				ch.flush()
			}
			ch.mutex.Unlock()
		}
	}
}

func (ch *batchingDedupedChannel[T, C]) flush() {
	batch := ch.buffer
	batchSlice := make([]T, 0, len(batch))
	for key, value := range batch {
		batchSlice = append(batchSlice, value)
		delete(ch.buffer, key)
	}
	ch.output <- batchSlice
}
