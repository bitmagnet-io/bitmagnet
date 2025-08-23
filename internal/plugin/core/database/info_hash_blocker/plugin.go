package info_hash_blocker

import (
	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/migrations"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	Blocker blocker.Blocker
}

var (
	Ref = database.Ref.MustSub("info_hash_blocker")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithDependencies[deps](
			migrations.Ref,
			postgres.Ref,
		),
		builder.WithEnabledByDefault[deps](),
		builder.WithFxOption[deps](
			fx.Provide(blocker.New),
		),
		builder.WithWorkerRegistryOption(
			func(deps deps) registry.Option {
				return registry.WithWorker(
					Ref.String(),
					deps.Blocker,
					worker.WithDependencies(
						migrations.Ref.String(),
						postgres.Ref.String(),
					),
				)
			},
		),
	)
)
