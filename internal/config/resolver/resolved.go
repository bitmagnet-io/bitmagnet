package resolver

import (
	"cmp"
	"maps"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/ref"
)

type Resolved struct {
	params map[string]*Param
}

func (r Resolved) Param(ref ref.Ref) (*Param, bool) {
	value, ok := r.params[ref.String()]

	return value, ok
}

func (r Resolved) Params() []*Param {
	params := slices.Collect(maps.Values(r.params))

	slices.SortFunc(params, func(a, b *Param) int {
		return cmp.Compare(a.Ref.String(), b.Ref.String())
	})

	return params
}
