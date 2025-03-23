package lazy

import "sync"

type Lazy[T any] interface {
	Get() (T, error)
	Decorate(func(T) (T, error))
	// IfInitialized calls the given function if the value has been initialized (useful for shutdown logic)
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

func (l *lazy[T]) Decorate(fn func(T) (T, error)) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	baseFn := l.fn
	l.fn = func() (T, error) {
		v, err := baseFn()
		if err != nil {
			return v, err
		}
		v, err = fn(v)
		return v, err
	}
}

func (l *lazy[T]) IfInitialized(fn func(T) error) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.done && l.err == nil {
		return fn(l.v)
	}
	return nil
}
