package ref

import (
	"cmp"
	"maps"
	"slices"
)

type Map[T any] struct {
	m map[string]entry[T]
}

func NewMap[T any]() Map[T] {
	return Map[T]{
		m: make(map[string]entry[T]),
	}
}

type Set = Map[struct{}]

func NewSet() Set {
	return NewMap[struct{}]()
}

type entry[T any] struct {
	Ref
	value T
}

func (m Map[T]) Set(ref Ref, value T) {
	m.m[ref.String()] = entry[T]{
		Ref:   ref,
		value: value,
	}
}

func (m Map[T]) SetAll(other Map[T]) {
	maps.Copy(m.m, other.m)
}

func (m Map[T]) Delete(ref Ref) {
	delete(m.m, ref.String())
}

func (m Map[T]) Get(ref Ref) T {
	value, _ := m.GetOK(ref)

	return value
}

func (m Map[T]) GetOK(ref Ref) (T, bool) {
	e, ok := m.m[ref.String()]

	return e.value, ok
}

func (m Map[T]) Has(ref Ref) bool {
	_, ok := m.GetOK(ref)

	return ok
}

func (m Map[T]) Refs() []Ref {
	refs := make([]Ref, 0, len(m.m))

	for _, e := range m.m {
		refs = append(refs, e.Ref)
	}

	slices.SortFunc(refs, func(a, b Ref) int {
		return cmp.Compare(a.String(), b.String())
	})

	return refs
}

func (m Map[T]) Values() []T {
	values := make([]T, 0, len(m.m))

	for _, ref := range m.Refs() {
		values = append(values, m.m[ref.String()].value)
	}

	return values
}

func (m Map[T]) StringMap() map[string]T {
	strMap := make(map[string]T, len(m.m))

	for key, entry := range m.m {
		strMap[key] = entry.value
	}

	return strMap
}

func (m Map[T]) Len() int {
	return len(m.m)
}
