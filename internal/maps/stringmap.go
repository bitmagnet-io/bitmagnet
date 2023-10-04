package maps

import (
	"github.com/facette/natsort"
	"sort"
)

type StringMap[T interface{}] map[string]T

type StringMapEntry[T interface{}] MapEntry[string, T]

// OrderedEntries returns the entries of the map in naturally sorted order by key.
func (m StringMap[T]) OrderedEntries() []StringMapEntry[T] {
	entries := make([]StringMapEntry[T], 0, len(m))
	for k, v := range m {
		entries = append(entries, StringMapEntry[T]{
			Key:   k,
			Value: v,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return natsort.Compare(entries[i].Key, entries[j].Key)
	})
	return entries
}

func (m StringMap[T]) WithValue(key string, value T) StringMap[T] {
	m[key] = value
	return m
}
