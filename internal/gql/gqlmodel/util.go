package gqlmodel

func nilToZero[T any](ptr *T) T {
	var result T
	if ptr != nil {
		result = *ptr
	}

	return result
}
