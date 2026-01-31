package deduplicator

import (
	"sync"
	"time"
)

type Deduplicator[T comparable] struct {
	mutex    sync.Mutex
	items    map[T]itemInfo
	keys     map[uint64]T
	minIndex uint64
	maxIndex uint64
	maxSize  int
	ttl      time.Duration
}

type itemInfo struct {
	addedAt time.Time
	index   uint64
}

func New[T comparable](maxSize int, ttl time.Duration) *Deduplicator[T] {
	return &Deduplicator[T]{
		items:   make(map[T]itemInfo, maxSize),
		keys:    make(map[uint64]T, maxSize),
		maxSize: max(1, maxSize),
		ttl:     ttl,
	}
}

func (d *Deduplicator[T]) Add(item T) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	info, exists := d.items[item]
	if exists {
		if time.Since(info.addedAt) < d.ttl {
			return false // Item is still valid
		}

		delete(d.items, item)
		delete(d.keys, info.index)

		if d.minIndex == info.index {
			d.minIndex++
		}
	}

	if len(d.items) >= d.maxSize {
		// Remove the oldest item
		for key, exists := d.keys[d.minIndex]; !exists; d.minIndex++ {
			delete(d.items, key)
			delete(d.keys, d.minIndex)
		}
	}

	d.items[item] = itemInfo{
		addedAt: time.Now(),
		index:   d.maxIndex,
	}
	d.keys[d.maxIndex] = item

	d.maxIndex++

	return true
}
