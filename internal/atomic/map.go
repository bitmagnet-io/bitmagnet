package atomic

type mappedReader[S any, T any] struct {
	source Reader[S]
	mapper func(S) T
}

func (mr mappedReader[S, T]) Get() T {
	return mr.mapper(mr.source.Get())
}

func (mr mappedReader[S, T]) Subscribe(fn func(T)) func() T {
	unsubscribe := mr.source.Subscribe(func(s S) {
		fn(mr.mapper(s))
	})

	return func() T {
		return mr.mapper(unsubscribe())
	}
}

func MapReader[S any, T any](source Reader[S], mapper func(S) T) Reader[T] {
	return mappedReader[S, T]{
		source: source,
		mapper: mapper,
	}
}

type intish interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func MapIntish[S intish, T intish](source Reader[S]) Reader[T] {
	return MapReader(source, func(s S) T {
		return T(s)
	})
}
