package atomic

type Reader[T any] interface {
	Get() T
	Subscribe(fn func(T)) func() T
}
