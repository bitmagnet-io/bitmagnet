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
