package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
}

var (
	Ref = database.Ref.MustSub("search")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides search functionality via the Postgres database"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithDependencies[deps](
			postgres.Ref,
		),
		builder.WithFxOption[deps](
			fx.Provide(search.New),
		),
	)
)
