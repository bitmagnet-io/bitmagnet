package lookup

type Chain []Lookup

func (l Chain) Lookup(path []string) (Result, bool, error) {
	var result *Result

	for i := len(l) - 1; i >= 0; i-- {
		resolver := l[i]

		thisResult, found, err := resolver.Lookup(path)
		if err != nil {
			return Result{}, false, err
		}

		if found {
			thisResult.prev = result
			result = &thisResult
		}
	}

	if result == nil {
		return Result{}, false, nil
	}

	return *result, true, nil
}
