// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lru

import (
	"sync"
	"time"
)

// EvictCallback is used to get a callback when a cache entry is evicted
type EvictCallback[K comparable, V any] func(key K, value V)

// LRU implements a thread-safe LRU with expirable entries.
type LRU[K comparable, V any] struct {
	size      int
	evictList *lruList[K, V]
	items     map[K]*entry[K, V]
	onEvict   EvictCallback[K, V]

	// expirable options
	mu  sync.Mutex
	ttl time.Duration

	chStopped chan struct{}
	stopOnce  sync.Once
	stopped   bool

	// buckets for expiration
	buckets []bucket[K, V]
	// uint8 because it's number between 0 and numBuckets
	nextCleanupBucket uint8
}

// bucket is a container for holding entries to be expired
type bucket[K comparable, V any] struct {
	entries     map[K]*entry[K, V]
	newestEntry time.Time
}

// because of uint8 usage for nextCleanupBucket, should not exceed 256.
// casting it as uint8 explicitly requires type conversions in multiple places
const numBuckets = 100

// New returns a new thread-safe cache with expirable entries.
//
// Size parameter set to 0 makes cache of unlimited size, e.g. turns LRU mechanism off.
//
// Providing 0 TTL turns expiring off.
//
// Delete expired entries every 1/100th of ttl value. Goroutine which deletes expired entries runs indefinitely.
func New[K comparable, V any](size int, onEvict EvictCallback[K, V], ttl time.Duration) *LRU[K, V] {
	if size < 0 {
		size = 0
	}

	res := LRU[K, V]{
		ttl:       ttl,
		size:      size,
		evictList: newList[K, V](),
		items:     make(map[K]*entry[K, V]),
		onEvict:   onEvict,
		chStopped: make(chan struct{}),
	}

	// initialize the buckets
	res.buckets = make([]bucket[K, V], numBuckets)
	for i := range numBuckets {
		res.buckets[i] = bucket[K, V]{entries: make(map[K]*entry[K, V])}
	}

	go func() {
		for {
			res.mu.Lock()
			ttl := res.ttl
			res.mu.Unlock()

			select {
			case <-res.chStopped:
				res.mu.Lock()
				res.stopped = true
				res.purge()
				res.mu.Unlock()

				return
			case <-time.After(ttl / numBuckets):
				res.mu.Lock()
				res.deleteExpired()
				res.mu.Unlock()
			}
		}
	}()

	return &res
}

// Purge clears the cache completely.
// onEvict is called for each evicted key.
func (c *LRU[K, V]) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped {
		return
	}

	c.purge()

	c.evictList.init()
}

func (c *LRU[K, V]) purge() {
	for k, v := range c.items {
		if c.onEvict != nil {
			c.onEvict(k, v.Value)
		}

		delete(c.items, k)
	}

	for _, b := range c.buckets {
		for _, ent := range b.entries {
			delete(b.entries, ent.Key)
		}
	}
}

// Add adds a value to the cache. Returns true if an eviction occurred.
// Returns false if there was no eviction: the item was already in the cache,
// the size was not exceeded or the cache was stopped.
func (c *LRU[K, V]) Add(key K, value V) (evicted bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped {
		return false
	}

	now := time.Now()

	// Check for existing item
	if ent, ok := c.items[key]; ok {
		c.evictList.moveToFront(ent)
		c.removeFromBucket(ent) // remove the entry from its current bucket as expiresAt is renewed
		ent.Value = value
		ent.AddedAt = now
		c.addToBucket(ent)

		return false
	}

	// Add new item
	ent := c.evictList.pushFrontExpirable(key, value, now.Add(c.ttl))
	c.items[key] = ent
	c.addToBucket(ent) // adds the entry to the appropriate bucket and sets entry.expireBucket

	evict := c.size > 0 && c.evictList.length() > c.size
	// Verify size not exceeded
	if evict {
		c.removeOldest()
	}

	return evict
}

// Get looks up a key's value from the cache.
func (c *LRU[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped {
		return
	}

	var ent *entry[K, V]
	if ent, ok = c.items[key]; ok {
		// Expired item check
		if time.Since(ent.AddedAt) > c.ttl {
			return value, false
		}

		c.evictList.moveToFront(ent)

		return ent.Value, true
	}

	return
}

// Contains checks if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *LRU[K, V]) Contains(key K) (ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped {
		return
	}

	_, ok = c.items[key]

	return ok
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *LRU[K, V]) Peek(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped {
		return
	}

	var ent *entry[K, V]
	if ent, ok = c.items[key]; ok {
		// Expired item check
		if time.Since(ent.AddedAt) > c.ttl {
			return value, false
		}

		return ent.Value, true
	}

	return
}

// Remove removes the provided key from the cache, returning if the
// key was contained.
func (c *LRU[K, V]) Remove(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped {
		return false
	}

	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
		return true
	}

	return false
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRU[K, V]) RemoveOldest() (key K, value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped {
		return
	}

	if ent := c.evictList.back(); ent != nil {
		c.removeElement(ent)
		return ent.Key, ent.Value, true
	}

	return
}

// GetOldest returns the oldest entry
func (c *LRU[K, V]) GetOldest() (key K, value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped {
		return
	}

	if ent := c.evictList.back(); ent != nil {
		return ent.Key, ent.Value, true
	}

	return
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
// Expired entries are filtered out.
func (c *LRU[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped {
		return nil
	}

	keys := make([]K, 0, len(c.items))

	now := time.Now()
	for ent := c.evictList.back(); ent != nil; ent = ent.prevEntry() {
		if now.After(ent.AddedAt.Add(c.ttl)) {
			continue
		}

		keys = append(keys, ent.Key)
	}

	return keys
}

// Values returns a slice of the values in the cache, from oldest to newest.
// Expired entries are filtered out.
func (c *LRU[K, V]) Values() []V {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped {
		return nil
	}

	values := make([]V, 0, len(c.items))

	now := time.Now()
	for ent := c.evictList.back(); ent != nil; ent = ent.prevEntry() {
		if now.After(ent.AddedAt.Add(c.ttl)) {
			continue
		}

		values = append(values, ent.Value)
	}

	return values
}

// Len returns the number of items in the cache.
func (c *LRU[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped {
		return 0
	}

	return c.evictList.length()
}

// Resize changes the cache size. Size of 0 means unlimited.
func (c *LRU[K, V]) Resize(size int) (evicted int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped {
		return
	}

	if size <= 0 {
		c.size = 0
		return 0
	}

	diff := c.evictList.length() - size
	if diff < 0 {
		diff = 0
	}

	for range diff {
		c.removeOldest()
	}

	c.size = size

	return diff
}

func (c *LRU[K, V]) SetTTL(ttl time.Duration) (evicted int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ttl = ttl

	return c.deleteExpired()
}

func (c *LRU[K, V]) Close() {
	c.stopOnce.Do(func() {
		close(c.chStopped)
	})
}

// removeOldest removes the oldest item from the cache. Has to be called with lock!
func (c *LRU[K, V]) removeOldest() {
	if ent := c.evictList.back(); ent != nil {
		c.removeElement(ent)
	}
}

// removeElement is used to remove a given list element from the cache. Has to be called with lock!
func (c *LRU[K, V]) removeElement(e *entry[K, V]) {
	c.evictList.remove(e)
	delete(c.items, e.Key)
	c.removeFromBucket(e)

	if c.onEvict != nil {
		c.onEvict(e.Key, e.Value)
	}
}

// deleteExpired deletes expired records.
func (c *LRU[K, V]) deleteExpired() int {
	removed := 0

	for {
		if c.stopped {
			return removed
		}

		bucketIdx := c.nextCleanupBucket
		timeToExpire := time.Until(c.buckets[bucketIdx].newestEntry.Add(c.ttl))

		if timeToExpire > 0 {
			return removed
		}

		shouldStop := true

		for _, ent := range c.buckets[bucketIdx].entries {
			c.removeElement(ent)

			removed++
			shouldStop = false
		}

		c.nextCleanupBucket = (c.nextCleanupBucket + 1) % numBuckets

		if shouldStop {
			return removed
		}
	}
}

// addToBucket adds entry to expire bucket so that it will be cleaned up when the time comes. Has to be called with
// lock!
func (c *LRU[K, V]) addToBucket(e *entry[K, V]) {
	bucketID := (numBuckets + c.nextCleanupBucket - 1) % numBuckets
	e.ExpireBucket = bucketID

	c.buckets[bucketID].entries[e.Key] = e
	if c.buckets[bucketID].newestEntry.Before(e.AddedAt) {
		c.buckets[bucketID].newestEntry = e.AddedAt
	}
}

// removeFromBucket removes the entry from its corresponding bucket. Has to be called with lock!
func (c *LRU[K, V]) removeFromBucket(e *entry[K, V]) {
	delete(c.buckets[e.ExpireBucket].entries, e.Key)
}

// Cap returns the capacity of the cache
func (c *LRU[K, V]) Cap() int {
	return c.size
}
