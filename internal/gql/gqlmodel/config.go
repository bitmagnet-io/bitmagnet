package gqlmodel

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/config/manager"
	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type ConfigQuery struct {
	I18n    *i18n.Bundle
	Manager *manager.Manager
}

func (q ConfigQuery) Params(ctx context.Context) []gen.ConfigParam {
	localizer := httpserver.NewLocalizerFromContext(ctx, q.I18n)

	return slice.Map(q.Manager.Params(), func(param *resolver.Param) gen.ConfigParam {
		return transformConfigParam(param, localizer)
	})
}

func (q ConfigQuery) Pending() bool {
	return q.Manager.HasPending()
}

func transformConfigParam(param *resolver.Param, localizer *i18n.Localizer) gen.ConfigParam {
	var description *string
	if localized, _ := localizer.LocalizeMessage(&i18n.Message{
		ID: param.Ref.String(),
	}); localized != "" {
		description = &localized
	} else {
		if strDesc := param.Description(); strDesc != "" {
			description = &strDesc
		}
	}

	defaultValue, err := param.EncodeYAMLAny(param.NewDefaultAny())
	if err != nil {
		panic(err)
	}

	return gen.ConfigParam{
		Ref:         param.Ref,
		Plugin:      param.Plugin,
		Description: description,
		Value:       json_schema.JSONValue(param.ValueYAML()),
		Source:      param.Source(),
		Default:     json_schema.JSONValue(defaultValue),
		Dynamic:     param.IsDynamic(),
		Pending:     param.IsPending(),
		JSONSchema:  param.JSONSchema(),
	}
}

type ConfigMutation struct {
	I18n    *i18n.Bundle
	Manager *manager.Manager
}

func (m ConfigMutation) Save(ctx context.Context, ref ref.Ref, value json_schema.JSONValue) (*gen.ConfigParam, error) {
	param, err := m.Manager.Save(ref, value.Value)
	if err != nil {
		return nil, err
	}

	result := transformConfigParam(param, httpserver.NewLocalizerFromContext(ctx, m.I18n))

	return &result, nil
}

func (m ConfigMutation) Delete(ctx context.Context, ref ref.Ref) (*gen.ConfigParam, error) {
	param, err := m.Manager.Delete(ref)
	if err != nil {
		return nil, err
	}

	result := transformConfigParam(param, httpserver.NewLocalizerFromContext(ctx, m.I18n))

	return &result, nil
}
