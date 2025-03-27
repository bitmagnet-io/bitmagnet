package maps

// InsertMap is a map that preserves initial insertion order.
type InsertMap[K comparable, V any] struct {
	keyValues map[K]V
	keys      []K
}

func NewInsertMap[K comparable, V any](entries ...MapEntry[K, V]) InsertMap[K, V] {
	m := InsertMap[K, V]{
		keys:      make([]K, 0, len(entries)),
		keyValues: make(map[K]V, len(entries)),
	}
	m.SetEntries(entries...)

	return m
}

func (m InsertMap[K, V]) Len() int {
	return len(m.keys)
}

func (m InsertMap[K, V]) Keys() []K {
	return m.keys
}

func (m InsertMap[K, V]) Values() []V {
	values := make([]V, 0, len(m.keys))
	for _, k := range m.keys {
		values = append(values, m.keyValues[k])
	}

	return values
}

func (m InsertMap[K, V]) Entries() []MapEntry[K, V] {
	values := make([]MapEntry[K, V], 0, len(m.keys))
	for _, k := range m.keys {
		values = append(values, MapEntry[K, V]{
			Key:   k,
			Value: m.keyValues[k],
		})
	}

	return values
}

func (m *InsertMap[K, V]) Set(key K, value V) {
	if _, ok := m.keyValues[key]; !ok {
		m.keys = append(m.keys, key)
	}

	m.keyValues[key] = value
}

func (m *InsertMap[K, V]) SetKey(key K) {
	var value V

	m.Set(key, value)
}

func (m *InsertMap[K, V]) SetEntries(entries ...MapEntry[K, V]) {
	for _, e := range entries {
		m.Set(e.Key, e.Value)
	}
}

func (m InsertMap[K, V]) Get(key K) (V, bool) {
	v, ok := m.keyValues[key]
	return v, ok
}

func (m InsertMap[K, V]) Has(key K) bool {
	_, ok := m.keyValues[key]
	return ok
}

func (m InsertMap[K, V]) Copy() InsertMap[K, V] {
	return NewInsertMap[K, V](m.Entries()...)
}
