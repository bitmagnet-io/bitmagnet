package ref

import (
	"cmp"
	"maps"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type Map[T any] struct {
	m map[string]Entry[T]
}

func NewMap[T any]() Map[T] {
	return Map[T]{
		m: make(map[string]Entry[T]),
	}
}

type Set = Map[struct{}]

func NewSet() Set {
	return NewMap[struct{}]()
}

type Entry[T any] struct {
	Ref
	Value T
}

func (m Map[T]) Set(ref Ref, value T) {
	m.m[ref.String()] = Entry[T]{
		Ref:   ref,
		Value: value,
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

	return e.Value, ok
}

func (m Map[T]) Entry(ref Ref) Entry[T] {
	e, _ := m.EntryOK(ref)

	return e
}

func (m Map[T]) EntryOK(ref Ref) (Entry[T], bool) {
	e, ok := m.m[ref.String()]

	return e, ok
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
		values = append(values, m.m[ref.String()].Value)
	}

	return values
}

func (m Map[T]) Entries() []Entry[T] {
	return slice.Map(m.Refs(), func(ref Ref) Entry[T] {
		return m.m[ref.String()]
	})
}

func (m Map[T]) StringMap() map[string]T {
	strMap := make(map[string]T, len(m.m))

	for key, entry := range m.m {
		strMap[key] = entry.Value
	}

	return strMap
}

func (m Map[T]) Len() int {
	return len(m.m)
}

func (m Map[T]) Clone() Map[T] {
	clone := NewMap[T]()
	clone.SetAll(m)

	return clone
}
