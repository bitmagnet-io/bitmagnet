package gqlmodel

import (
	"context"
	"slices"

	"github.com/99designs/gqlgen/graphql"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
)

type PluginQuery struct {
	Plugins ref.Map[plugin.Instance]
}

func (q PluginQuery) List(ctx context.Context) ([]gen.PluginInfo, error) {
	hasLocalizedContent := slices.Contains(graphql.CollectAllFields(ctx), "localizedContent")

	var acceptLanguage []string
	if hasLocalizedContent {
		acceptLanguage = httpserver.AcceptLanguageFromContext(ctx)
	}

	return slice.MapErr(q.Plugins.Values(), func(instance plugin.Instance) (gen.PluginInfo, error) {
		var localizedContent plugin.LocalizedContent

		if hasLocalizedContent {
			var err error

			localizedContent, err = instance.LocalizedContent(ctx, acceptLanguage...)
			if err != nil {
				return gen.PluginInfo{}, err
			}
		}

		return transformPluginInfo(instance, localizedContent), nil
	})
}

func transformPluginInfo(
	instance plugin.Instance,
	localizedContent plugin.LocalizedContent,
) gen.PluginInfo {
	var description *string
	if localizedContent.Description != "" {
		description = &localizedContent.Description
	}

	return gen.PluginInfo{
		Ref:        instance.Ref(),
		Enabled:    instance.Enabled(),
		DependsOn:  instance.DependsOn(),
		RequiredBy: instance.RequiredBy(),
		LocalizedContent: gen.PluginLocalizedContent{
			Description: description,
			ConfigParams: slice.Map(
				localizedContent.ConfigParams,
				func(param plugin.LocalizedConfigParam) gen.ConfigParamLocalizedContent {
					return gen.ConfigParamLocalizedContent{
						Ref:         param.Ref,
						Description: param.Description,
					}
				},
			),
		},
	}
}
