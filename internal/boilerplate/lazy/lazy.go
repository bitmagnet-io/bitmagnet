package lazy

import "sync"

type Lazy[T any] interface {
	Get() (T, error)
	IfInitialized(func(T) error) error
}

func New[T any](fn func() (T, error)) Lazy[T] {
	return &lazy[T]{fn: fn}
}

type lazy[T any] struct {
	fn   func() (T, error)
	mtx  sync.Mutex
	v    T
	err  error
	done bool
}

func (l *lazy[T]) Get() (T, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if !l.done {
		l.v, l.err = l.fn()
		l.done = true
	}
	return l.v, l.err
}

func (l *lazy[T]) IfInitialized(fn func(T) error) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.done {
		return fn(l.v)
	}
	return nil
}
