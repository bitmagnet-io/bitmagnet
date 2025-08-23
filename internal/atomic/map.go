package atomic

// Map maps a *Value[From] to a *Value[To] using the provided transformer.
// It returns the mapped Value along with an unsubscribe function that separates it from the input value.
func Map[From any, To any](
	value *Value[From],
	transform func(From) To,
) (*Value[To], func() From) {
	mapped := &Value[To]{}
	unsubscribe := value.Subscribe(func(from From) {
		mapped.Set(transform(from))
	})

	return mapped, unsubscribe
}
