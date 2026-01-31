package search

import (
	i18n_embed "github.com/bitmagnet-io/bitmagnet/i18n"
	db_search "github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/search"
	"github.com/bitmagnet-io/bitmagnet/internal/search/adapter/multi"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/pkg/i18n"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
}

var (
	Ref = ref.Root.MustSub("search")

	facetRef = Ref.MustSub("facet")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides search functionality via the Postgres database"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithDependencies[deps](
			postgres.Ref,
		),
		builder.WithOptions(
			slice.Map(search.FacetValues(), func(facet search.Facet) builder.Option[deps] {
				return builder.WithI18nMessage[deps](
					facetRef.MustSub(facet.String()),
					"label for search facet: "+facet.String(),
					i18n.WithOther(facet.Label()),
				)
			})...,
		),
		builder.WithFxOption[deps](
			fx.Provide(func() search.FacetLocalizer {
				return search.NewFacetLocalizer(facetRef, i18n_embed.Bundle)
			}),
			fx.Provide(
				db_search.New,
				fx.Annotate(
					func(adapters [][]multi.Index, s db_search.Search) multi.Search {
						return multi.New(append(
							[]multi.Index{
								{
									Ref:     Ref.MustSub("postgres"),
									Name:    "Postgres",
									Adapter: db_search.Adapter{Search: s},
								},
							},
							slice.Flatten(adapters...)...,
						)...)
					},
					fx.ParamTags(`group:"search_adapters"`),
					fx.As(new(search.Search)),
					fx.As(new(search.TorrentContent)),
					fx.As(new(search.TorrentFiles)),
					fx.As(new(search.Content)),
					fx.As(fx.Self()),
				),
			),
		),
	)
)
