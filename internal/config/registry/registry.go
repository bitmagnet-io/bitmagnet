package registry

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
)

type (
	// Registry struct {
	// 	params map[string]Param
	// }

	Param struct {
		ref.Ref
		Plugin ref.Ref
		param.Untyped
	}
)

// func New(options []Option) (Registry, error) {
// 	r := Registry{
// 		params: make(map[string]Param),
// 	}

// 	var err error

// 	for _, opt := range options {
// 		err = opt(&r)

// 		if err != nil {
// 			return r, err
// 		}
// 	}

// 	return r, nil
// }

// func (r Registry) All() []Param {
// 	result := slices.Collect(maps.Values(r.params))

// 	slices.SortFunc(result, func(a, b Param) int {
// 		return cmp.Compare(a.Ref.String(), b.Ref.String())
// 	})

// 	return result
// }

// func (r Registry) Param(ref ref.Ref) (Param, bool) {
// 	p, ok := r.params[ref.String()]
// 	return Param{
// 		Ref:     ref,
// 		Untyped: p,
// 	}, ok
// }

// type Option func(*Registry) error

// func WithParam(ref ref.Ref, param param.Untyped) Option {
// 	return func(r *Registry) error {
// 		str := ref.String()
// 		if _, ok := r.params[str]; ok {
// 			return fmt.Errorf("parameter already exists: %s", str)
// 		}
// 		r.params[str] = Param{
// 			Ref:     ref,
// 			Untyped: param,
// 		}

// 		return nil
// 	}
// }
