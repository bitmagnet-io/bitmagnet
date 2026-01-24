package slice

func Map[T any, R any](t []T, mapFunc func(T) R) []R {
	r := make([]R, len(t))

	for i, e := range t {
		r[i] = mapFunc(e)
	}

	return r
}

func MapErr[T any, R any](t []T, mapFunc func(T) (R, error)) ([]R, error) {
	var err error

	r := make([]R, len(t))

	for i, e := range t {
		r[i], err = mapFunc(e)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func FlatMap[T any, R any](t []T, mapFunc func(T) []R) []R {
	//nolint:prealloc
	var result []R

	for _, e := range t {
		result = append(result, mapFunc(e)...)
	}

	return result
}

func MapWithArg[I any, O any, A any](t []I, arg A, mapFunc func(A, I) O) []O {
	return Map(t, func(e I) O {
		return mapFunc(arg, e)
	})
}

func Filter[T any](t []T, predicate func(T) bool) []T {
	var result []T

	for _, e := range t {
		if predicate(e) {
			result = append(result, e)
		}
	}

	return result
}

func Some[T any](t []T, predicate func(T) bool) bool {
	for _, e := range t {
		if predicate(e) {
			return true
		}
	}

	return false
}

func Group[T any, K comparable](s []T, keyFunc func(T) K) map[K][]T {
	m := map[K][]T{}

	for _, item := range s {
		k := keyFunc(item)
		m[k] = append(m[k], item)
	}

	return m
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
