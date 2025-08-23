package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
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

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithDependencies[deps](
			config.Ref,
			search.Ref,
			tmdb_compat.Ref,
		),
		// builder.WithDefaultConfig[deps](classifier.NewDefaultConfig()),
		builder.WithConfigParam[deps](Ref.MustSub("workflow"), classifier.ParamWorkflow),
		builder.WithConfigParam[deps](Ref.MustSub("keywords"), classifier.ParamKeywords),
		builder.WithConfigParam[deps](Ref.MustSub("extensions"), classifier.ParamExtensions),
		builder.WithConfigParam[deps](Ref.MustSub("flags"), classifier.ParamFlags),
		builder.WithConfigParam[deps](Ref.MustSub("delete_xxx"), classifier.ParamDeleteXXX),
		builder.WithConfigParam[deps](Ref.MustSub("concurrency"), classifier.ParamConcurrency),
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
