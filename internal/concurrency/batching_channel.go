package concurrency

import (
  "sync"
  "time"
)

type BatchingChannel[T any] interface {
	In() chan<- T
	Out() <-chan []T
}

type batchingChannel[T any] struct {
	mutex        sync.Mutex
	input        chan T
	output       chan []T
	buffer       []T
	maxBatchSize int
	maxWaitTime  time.Duration
	ticker       *time.Ticker
}

func NewBatchingChannel[T any](input chan T, maxBatchSize int, maxWaitTime time.Duration) BatchingChannel[T] {
	ch := &batchingChannel[T]{
		input:        input,
		output:       make(chan []T),
		maxBatchSize: maxBatchSize,
		maxWaitTime:  maxWaitTime,
		ticker:       time.NewTicker(maxWaitTime),
	}
	go ch.batch()
	return ch
}

func (ch *batchingChannel[T]) In() chan<- T {
	return ch.input
}

func (ch *batchingChannel[T]) Out() <-chan []T {
	return ch.output
}

func (ch *batchingChannel[T]) batch() {
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
			ch.buffer = append(ch.buffer, next)
			if len(ch.buffer) >= ch.maxBatchSize {
				ch.flush()
			}
			ch.mutex.Unlock()
		case <-ch.ticker.C:
			ch.mutex.Lock()
			if len(ch.buffer) >= 0 {
				ch.flush()
			}
			ch.mutex.Unlock()
		}
	}
}

func (ch *batchingChannel[T]) flush() {
	ch.ticker.Stop()
	batch := ch.buffer
	ch.buffer = nil
	ch.output <- batch
	ch.ticker.Reset(ch.maxWaitTime)
}
