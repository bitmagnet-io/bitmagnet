package resolver

import (
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
)

type Resolved struct {
	params ref.Map[*Param]
}

func (r Resolved) Param(ref ref.Ref) (*Param, bool) {
	return r.params.GetOK(ref)
}

func (r Resolved) Params() []*Param {
	return r.params.Values()
}
