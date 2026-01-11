package gqlmodel

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
)

type PluginQuery struct {
	I18n  *i18n.Bundle
	Infos plugin.PluginInfos
}

func (q PluginQuery) List(ctx context.Context) []gen.PluginInfo {
	localizer := httpserver.NewLocalizerFromContext(ctx, q.I18n)

	return slice.Map(q.Infos, func(info plugin.PluginInfo) gen.PluginInfo {
		return transformPluginInfo(info, localizer)
	})
}

func transformPluginInfo(info plugin.PluginInfo, localizer *i18n.Localizer) gen.PluginInfo {
	var description *string
	if localized, _ := localizer.LocalizeMessage(&i18n.Message{
		ID: info.Ref.String(),
	}); localized != "" {
		description = &localized
	}

	return gen.PluginInfo{
		Ref:         info.Ref,
		Description: description,
		Enabled:     info.Enabled,
		DependsOn:   info.DependsOn,
		RequiredBy:  info.RequiredBy,
	}
}
