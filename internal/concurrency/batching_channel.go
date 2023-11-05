package concurrency

import (
	"time"
)

type BatchingChannel[T any] interface {
	In() chan<- T
	Out() <-chan []T
}

type batchingChannel[T any] struct {
	input        chan T
	output       chan []T
	buffer       []T
	maxBatchSize int
	maxWaitTime  time.Duration
	ticker       *time.Ticker
}

func NewBatchingChannel[T any](capacity int, maxBatchSize int, maxWaitTime time.Duration) BatchingChannel[T] {
	ch := &batchingChannel[T]{
		input:        make(chan T, capacity),
		output:       make(chan []T, 1),
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
			ch.buffer = append(ch.buffer, next)
			if len(ch.buffer) >= ch.maxBatchSize {
				ch.flush()
			}
		case <-ch.ticker.C:
			if len(ch.buffer) > 0 {
				ch.flush()
			}
		}
	}
}

func (ch *batchingChannel[T]) flush() {
	ch.ticker.Stop()
	batch := ch.buffer
	ch.buffer = nil
	ch.ticker.Reset(ch.maxWaitTime)
	ch.output <- batch
}
