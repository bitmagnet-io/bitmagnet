package maps

type MapEntry[K comparable, V interface{}] struct {
	Key   K
	Value V
}
