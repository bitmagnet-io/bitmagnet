package slice

import (
	"iter"
)

func Map[T any, R any](t []T, mapFunc func(T) R) []R {
	r := make([]R, len(t))

	for i, e := range t {
		r[i] = mapFunc(e)
	}

	return r
}

func MapWithArg[I any, O any, A any](t []I, arg A, mapFunc func(A, I) O) []O {
	return Map(t, func(e I) O {
		return mapFunc(arg, e)
	})
}

func Group[T any, K comparable](s []T, keyFunc func(T) K) map[K][]T {
	m := map[K][]T{}

	for _, item := range s {
		k := keyFunc(item)
		m[k] = append(m[k], item)
	}

	return m
}

func ToMap[T any, K comparable, V any](s []T, transformFunc func(T) (K, V)) map[K]V {
	m := make(map[K]V, len(s))

	for _, item := range s {
		k, v := transformFunc(item)
		m[k] = v
	}

	return m
}

func Insert[T any](slice []T, value T, index int) []T {
	return append(slice[:index], append([]T{value}, slice[index:]...)...)
}

func Remove[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func Unique[T comparable](list []T) []T {
	var result []T

	seen := make(map[T]struct{})

	for _, item := range list {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}

			result = append(result, item)
		}
	}

	return result
}

// CollectChunks collects chunks of n elements from the input sequence and return a Seq of chunks
func CollectChunks[T any](it iter.Seq[T], n int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		s := make([]T, 0, n)

		for x := range it {
			s = append(s, x)
			if len(s) >= n {
				if !yield(s) {
					return
				}

				s = make([]T, 0, n)
			}
		}

		if len(s) > 0 {
			yield(s)
		}
	}
}

// SeqFunc returns a Seq that iterates over the slice with the given mapping function
func SeqFunc[I, O any](s []I, f func(I) O) iter.Seq[O] {
	return func(yield func(O) bool) {
		for _, x := range s {
			if !yield(f(x)) {
				return
			}
		}
	}
}
