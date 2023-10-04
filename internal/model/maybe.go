package model

type Maybe[T interface{}] struct {
	Val   T
	Valid bool
}

func MaybeValid[T interface{}](v T) Maybe[T] {
	return Maybe[T]{Val: v, Valid: true}
}

func (n Maybe[T]) IsValid() bool {
	return n.Valid
}

func (n Maybe[T]) IsDefined() bool {
	return n.Valid
}

func (n Maybe[T]) Addr() *T {
	if !n.Valid {
		return nil
	}

	return &n.Val
}
