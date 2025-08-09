package persister

import (
	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	internal_database "github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/info_hash_blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/migrations"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	plugin_metrics "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/metrics"
	plugin_worker "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/worker"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	Config = persister.Config

	deps struct {
		fx.In
		Worker persister.Persister
	}
)

var (
	Ref = database.Ref.MustSub("persister")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithDependencies[Config, deps](
			info_hash_blocker.Ref,
			logging.Ref,
			plugin_metrics.Ref,
			migrations.Ref,
			plugin_worker.Ref,
			postgres.Ref,
		),
		builder.WithDefaultConfig[Config, deps](persister.NewDefaultConfig()),
		builder.WithFxOption[Config, deps](
			fx.Provide(
				fx.Annotate(
					func(
						config Config,
						daoProvider internal_database.DaoTransactionProvider,
						blockerBlocker blocker.Blocker,
						logger *zap.Logger,
						metrics *metrics.Registry,
					) persister.Persister {
						return persister.New(
							config,
							daoProvider,
							blockerBlocker,
							logger.Named(Ref.String()),
							metrics.MustNewComponent(Ref),
						)
					},
					fx.As(new(persister.Adder)),
					fx.As(new(persister.Persister)),
				),
			),
		),
		builder.WithWorkerRegistryOption(
			func(_ Config, deps deps) registry.Option {
				return registry.WithWorker(
					Ref.String(),
					deps.Worker,
					worker.WithDependencies(
						postgres.Ref.String(),
						info_hash_blocker.Ref.String(),
						migrations.Ref.String(),
					),
				)
			},
		),
	)
)
