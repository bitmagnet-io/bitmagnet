package concurrency

import (
	"context"
	"errors"
)

type BufferedDedupedChannel[T any, C comparable] interface {
	In(context.Context, T) error
	TryIn(T) error
	Out() <-chan T
	Len() int
	SpareCap() int
}

func NewBufferedDedupedChannel[T any, C comparable](
	capacity int,
	toComparable func(T) C,
) BufferedDedupedChannel[T, C] {
	return &bufferedDedupedChannel[T, C]{
		toComparable: toComparable,
		cap:          capacity,
		ch:           make(chan T, capacity),
		set:          &AtomicSet[C]{},
	}
}

type bufferedDedupedChannel[T any, C comparable] struct {
	toComparable func(T) C
	cap          int
	ch           chan T
	set          *AtomicSet[C]
}

var ErrDuplicateItem = errors.New("duplicate item")
var ErrChannelFull = errors.New("channel full")

func (c *bufferedDedupedChannel[T, C]) In(ctx context.Context, value T) error {
	cmp := c.toComparable(value)
	if !c.set.Add(cmp) {
		return ErrDuplicateItem
	}
	defer c.set.Remove(cmp)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.ch <- value:
		return nil
	}
}

func (c *bufferedDedupedChannel[T, C]) TryIn(value T) error {
	cmp := c.toComparable(value)
	if !c.set.Add(cmp) {
		return ErrDuplicateItem
	}
	defer c.set.Remove(cmp)
	select {
	case c.ch <- value:
		return nil
	default:
		return ErrChannelFull
	}
}

func (c *bufferedDedupedChannel[T, C]) Out() <-chan T {
	return c.ch
}

func (c *bufferedDedupedChannel[T, C]) Len() int {
	return c.set.Len()
}

func (c *bufferedDedupedChannel[T, C]) SpareCap() int {
	return c.cap - c.set.Len()
}
