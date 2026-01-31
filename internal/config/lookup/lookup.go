package lookup

type Result struct {
	source string
	value  any
	prev   *Result
}

func (r Result) Key() string {
	return r.source
}

func (r Result) Value() any {
	return r.value
}

func (r Result) Prev() (Result, bool) {
	if r.prev != nil {
		return *r.prev, true
	}

	return Result{}, false
}

func (r Result) Chain() []Result {
	var chain []Result

	for current := &r; current != nil && current.source != ""; current = current.prev {
		chain = append(chain, *current)
	}

	return chain
}

type Lookup interface {
	Lookup(path []string) (
		result Result,
		found bool,
		err error,
	)
}
