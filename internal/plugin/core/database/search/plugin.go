package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"go.uber.org/fx"
)

type (
	config struct{}
	deps   struct{}
)

var (
	Ref = database.Ref.MustSub("search")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[config, deps](),
		builder.WithDependencies[config, deps](
			postgres.Ref,
		),
		builder.WithFxOption[config, deps](
			fx.Provide(search.New),
		),
	)
)
