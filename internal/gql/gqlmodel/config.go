package gqlmodel

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/manager"
	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type ConfigQuery struct {
	Manager *manager.Manager
}

func (q ConfigQuery) Params() []gen.ConfigParam {
	return slice.Map(q.Manager.Params(), transformConfigParam)
}

func (q ConfigQuery) Pending() bool {
	return q.Manager.HasPending()
}

func transformConfigParam(param *resolver.Param) gen.ConfigParam {
	var doc *string
	if strDoc := param.Doc(); strDoc != "" {
		doc = &strDoc
	}
	return gen.ConfigParam{
		Ref:     param.Ref.String(),
		Doc:     doc,
		Value:   param.ValueString(),
		Source:  param.Source(),
		Default: param.Last().ValueString(),
		Dynamic: param.IsDynamic(),
		Pending: param.IsPending(),
	}
}

type ConfigMutation struct {
	Manager *manager.Manager
}

func (m ConfigMutation) Save(key string, value string) (*gen.ConfigParam, error) {
	ref, err := ref.Parse(key)
	if err != nil {
		return nil, err
	}

	param, err := m.Manager.Save(ref, value)
	if err != nil {
		return nil, err
	}

	result := transformConfigParam(param)

	return &result, nil
}

func (m ConfigMutation) Delete(key string) (*gen.ConfigParam, error) {
	ref, err := ref.Parse(key)
	if err != nil {
		return nil, err
	}

	param, err := m.Manager.Delete(ref)
	if err != nil {
		return nil, err
	}

	result := transformConfigParam(param)

	return &result, nil
}
