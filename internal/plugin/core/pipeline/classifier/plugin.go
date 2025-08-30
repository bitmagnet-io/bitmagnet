package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline"
	tmdb_compat "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/tmdb/compat"
	"go.uber.org/fx"
)

type deps = classifier.Params

var (
	Ref = pipeline.Ref.MustSub("classifier")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides the torrent classifier service"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithDependencies[deps](
			config.Ref,
			search.Ref,
			tmdb_compat.Ref,
		),
		// builder.WithDefaultConfig[deps](classifier.NewDefaultConfig()),
		builder.WithConfig[deps](Ref.MustSub("workflow"), classifier.ParamWorkflow),
		builder.WithConfig[deps](Ref.MustSub("keywords"), classifier.ParamKeywords),
		builder.WithConfig[deps](Ref.MustSub("extensions"), classifier.ParamExtensions),
		builder.WithConfig[deps](Ref.MustSub("flags"), classifier.ParamFlags),
		builder.WithConfig[deps](Ref.MustSub("delete_xxx"), classifier.ParamDeleteXXX),
		builder.WithConfig[deps](Ref.MustSub("concurrency"), classifier.ParamConcurrency),
		builder.WithFxOption[deps](
			fx.Provide(
				func(
					workflow classifier.Workflow,
					keywords classifier.Keywords,
					extensions classifier.Extensions,
					flags classifier.FlagValues,
					deleteXXX classifier.DeleteXXX,
					concurrency classifier.Concurrency,
				) classifier.Config {
					return classifier.Config{
						Workflow:    workflow,
						Keywords:    keywords,
						Extensions:  extensions,
						Flags:       flags,
						DeleteXXX:   deleteXXX,
						Concurrency: concurrency,
					}
				},
				classifier.New,
			),
		),
	)
)
