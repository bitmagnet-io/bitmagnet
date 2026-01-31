package batch

import (
	"context"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

type workerRunner[K comparable, V any] struct {
	*worker[K, V]
	ctx       context.Context
	cancel    context.CancelCauseFunc
	buffer    maps.InsertMap[K, value[V]]
	lastFlush time.Time
}

type value[V any] struct {
	value V
	added time.Time
}

func (wr *workerRunner[K, V]) run() {
	defer func() {
		wr.cancel(nil)
		wr.metrics.Reset()
	}()

	for {
		select {
		case <-wr.shutdown:
			_ = wr.addItemsAndCheckFlush(nil, true)
			return
		case <-time.After(wr.maxWait):
			_ = wr.addItemsAndCheckFlush(nil, true)
		case items, ok := <-wr.ch:
			if !ok {
				return
			}

			err := wr.addItemsAndCheckFlush(items, false)
			if err != nil {
				wr.cancel(err)
				return
			}
		}
	}
}

func (wr *workerRunner[K, V]) addItemsAndCheckFlush(items []V, forceFlush bool) error {
	now := time.Now()

	wr.mtx.Lock()

	initialSize := wr.buffer.Len()

	for _, item := range items {
		key := wr.keyer(item)

		if !wr.filter(key, item) {
			wr.metrics.IncrDiscarded()
			continue
		}

		var (
			v  value[V]
			ok bool
		)

		if v, ok = wr.buffer.Get(key); ok {
			v.value = wr.merger(v.value, item)
			wr.metrics.IncrDeduplicated()
		} else {
			v = value[V]{value: item, added: now}
		}

		wr.buffer.Set(key, v)
	}

	size := wr.buffer.Len()

	wr.metrics.IncrAdded(size - initialSize)

	shouldFlush := size > 0 && (forceFlush || size > wr.maxSize || time.Since(wr.lastFlush) > wr.maxWait)

	var values []value[V]

	if shouldFlush {
		values = wr.buffer.Values()
		wr.buffer.Clear()
	}

	wr.lastFlush = time.Now()

	wr.mtx.Unlock()

	if shouldFlush {
		now = time.Now()
		items = make([]V, 0, len(values))

		for _, v := range values {
			items = append(items, v.value)
			wr.metrics.IncrDequeued(now.Sub(v.added))
		}

		err := wr.flusher(wr.ctx, items)
		if err != nil {
			return fmt.Errorf("%w: %w: %w", Err, ErrFlush, err)
		}

		now = time.Now()

		for _, v := range values {
			wr.metrics.IncrFlushed(now.Sub(v.added))
		}
	}

	return nil
}

func (wr *workerRunner[K, V]) shutdowner(ctx context.Context) error {
	close(wr.shutdown)

	var err error

	if wr.quickShutdown {
		wr.cancel(runner.ErrShutdownRequested)
	} else {
		select {
		case <-ctx.Done():
			err = fmt.Errorf("%w: %w: %w", Err, ErrShutdown, ctx.Err())
		case <-wr.ctx.Done():
		}
	}

	wr.mtx.Lock()
	ch := wr.ch
	wr.ch = nil
	wr.shutdown = nil
	wr.mtx.Unlock()

	close(ch)

	return err
}
