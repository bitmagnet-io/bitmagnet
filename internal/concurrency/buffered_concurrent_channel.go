package concurrency

import (
	"context"
	"golang.org/x/sync/semaphore"
)

type BufferedConcurrentChannel[T any] interface {
	In() chan<- T
	Run(context.Context, func(T)) error
}

func NewBufferedConcurrentChannel[T any](capacity int, concurrency int) BufferedConcurrentChannel[T] {
	return bufferedConcurrentChannel[T]{
		ch:  make(chan T, capacity),
		sem: semaphore.NewWeighted(int64(concurrency)),
	}
}

type bufferedConcurrentChannel[T any] struct {
	ch  chan T
	sem *semaphore.Weighted
}

func (ch bufferedConcurrentChannel[T]) In() chan<- T {
	return ch.ch
}

func (ch bufferedConcurrentChannel[T]) Run(ctx context.Context, f func(T)) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case next := <-ch.ch:
			if err := ch.sem.Acquire(ctx, 1); err != nil {
				return err
			}
			go func() {
				defer ch.sem.Release(1)
				f(next)
			}()
		}
	}
}
