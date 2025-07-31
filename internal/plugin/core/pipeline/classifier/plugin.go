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

type (
	Config = classifier.Config
	deps   = classifier.Params
)

var (
	Ref = pipeline.Ref.MustSub("classifier")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithDependencies[Config, deps](
			config.Ref,
			search.Ref,
			tmdb_compat.Ref,
		),
		builder.WithDefaultConfig[Config, deps](classifier.NewDefaultConfig()),
		builder.WithFxOption[Config, deps](fx.Provide(classifier.New)),
	)
)
