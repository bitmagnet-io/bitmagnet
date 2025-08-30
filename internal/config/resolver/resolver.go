package resolver

import (
	"reflect"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/config"
	"github.com/bitmagnet-io/bitmagnet/internal/config/lookup"
	"github.com/bitmagnet-io/bitmagnet/internal/config/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"gopkg.in/yaml.v3"
)

type Resolver struct {
	lookup lookup.Lookup
	params []registry.Param
}

func New(lookup lookup.Lookup, params ...registry.Param) Resolver {
	return Resolver{
		lookup: lookup,
		params: params,
	}
}

func (r Resolver) Resolve() (Resolved, error) {
	params := ref.NewMap[*Param]()

	for _, param := range r.params {
		resolved, err := r.resolve(param)
		if err != nil {
			return Resolved{}, err
		}

		params.Set(param.Ref, resolved)
	}

	return Resolved{params: params}, nil
}

func (r Resolver) resolve(param registry.Param) (*Param, error) {
	lookup, _, err := r.lookup.Lookup(param.Ref.Path())
	if err != nil {
		return nil, err
	}

	res := newParam(param, config.SourceDefault, param.NewDefaultAny(), nil)

	lookupChain := lookup.Chain()

	for i := len(lookupChain) - 1; i >= 0; i-- {
		item := lookupChain[i]
		rawValue := item.Value()

		var (
			value any
			err   error
		)

		switch v := rawValue.(type) {
		case string:
			value, err = param.ParseAny(v)
		case yaml.Node:
			value, err = param.DecodeYAMLAny(v)
		}

		if err == nil {
			err = param.ValidateAny(value)
		}

		if err != nil {
			// If the error concerns an overridden config value,
			// exclude from the chain and continue; otherwise fail.
			if i > 0 {
				continue
			}
			return nil, err
		}

		prev := res.chain
		res = newParam(param, item.Key(), value, prev)
	}

	return res, nil
}

func newParam(param registry.Param, lookupKey string, value any, prev *chain) *Param {
	var initialValue any

	if reflected, ok := atomic.ReflectValue(reflect.ValueOf(value)); ok {
		initialValue = reflected.GetAny()
	} else {
		initialValue = value
	}

	return &Param{
		param: param,
		chain: &chain{
			source:       lookupKey,
			value:        value,
			initialValue: initialValue,
			prev:         prev,
		},
	}
}
