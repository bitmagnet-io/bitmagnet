package concurrency

import "sync"

type AtomicValue[T any] struct {
	mutex sync.RWMutex
	value T
}

func (a *AtomicValue[T]) Get() T {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.value
}

func (a *AtomicValue[T]) Set(value T) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.value = value
}

func (a *AtomicValue[T]) Update(fn func(T) T) T {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.value = fn(a.value)
	return a.value
}

type AtomicCounter struct {
	AtomicValue[int]
}

func (c *AtomicCounter) Inc(n int) int {
	return c.Update(func(v int) int {
		return v + n
	})
}

type AtomicSet[T comparable] struct {
	AtomicValue[map[T]struct{}]
}

func (s *AtomicSet[T]) Add(value T) bool {
	added := false
	s.Update(func(v map[T]struct{}) map[T]struct{} {
		if _, ok := v[value]; ok {
			return v
		}
		added = true
		newV := make(map[T]struct{}, len(v)+1)
		for k := range v {
			newV[k] = struct{}{}
		}
		newV[value] = struct{}{}
		return newV
	})
	return added
}

func (s *AtomicSet[T]) AddUpTo(value T, capacity int) bool {
	added := false
	s.Update(func(v map[T]struct{}) map[T]struct{} {
		if len(v) > capacity {
			return v
		}
		if _, ok := v[value]; ok {
			return v
		}
		added = true
		newV := make(map[T]struct{}, len(v)+1)
		for k := range v {
			newV[k] = struct{}{}
		}
		newV[value] = struct{}{}
		return newV
	})
	return added
}

func (s *AtomicSet[T]) Remove(value T) bool {
	removed := false
	s.Update(func(v map[T]struct{}) map[T]struct{} {
		if _, ok := v[value]; !ok {
			return v
		}
		removed = true
		newV := make(map[T]struct{}, len(v)-1)
		for k := range v {
			if k != value {
				newV[k] = struct{}{}
			}
		}
		return newV
	})
	return removed
}

func (s *AtomicSet[T]) Has(value T) bool {
	v := s.Get()
	if v == nil {
		return false
	}
	_, ok := v[value]
	return ok
}

func (s *AtomicSet[T]) Len() int {
	v := s.Get()
	return len(v)
}
