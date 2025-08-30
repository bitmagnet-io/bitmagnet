package atomic

import (
	"fmt"
	"sync"
)

type Value[T any] struct {
	mutex           sync.RWMutex
	value           T
	subscribers     map[uint64]func(T)
	subscriberIndex uint64
}

type GetAny interface {
	GetAny() any
}

type SetAny interface {
	SetAny(value any) error
}

type GetSetAny interface {
	GetAny
	SetAny
}

var _ GetAny = (*Value[any])(nil)

func NewValue[T any](initialValue T) *Value[T] {
	return &Value[T]{
		value: initialValue,
	}
}

func (v *Value[T]) Get() T {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	return v.value
}

func (v *Value[T]) GetAny() any {
	return v.Get()
}

func (v *Value[T]) Set(value T) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.set(value)
}

func (v *Value[T]) SetAny(value any) error {
	if atomicValue, ok := value.(*Value[T]); ok {
		value = atomicValue.Get()
	}

	if typedValue, ok := value.(T); ok {
		v.Set(typedValue)
		return nil
	}

	return fmt.Errorf("expected value of type %T, got %T", v.value, value)
}

func (v *Value[T]) Update(fn func(T) T) T {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.set(fn(v.value))

	return v.value
}

func (v *Value[T]) set(value T) {
	v.value = value

	for _, subscriber := range v.subscribers {
		subscriber(value)
	}
}

func (v *Value[T]) Subscribe(fn func(T)) func() T {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	index := v.subscriberIndex

	if index == 0 {
		v.subscribers = make(map[uint64]func(T))
	}

	v.subscribers[index] = fn
	v.subscriberIndex++

	fn(v.value)

	return func() T {
		v.mutex.Lock()
		defer v.mutex.Unlock()

		delete(v.subscribers, index)

		return v.value
	}
}
