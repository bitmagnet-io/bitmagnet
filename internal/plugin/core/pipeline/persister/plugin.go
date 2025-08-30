package persister

import (
	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	internal_database "github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/info_hash_blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/migrator"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	plugin_metrics "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/metrics"
	plugin_worker "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/worker"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type deps struct {
	fx.In
	Worker persister.Persister
}

var (
	Ref = database.Ref.MustSub("persister")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Runs the worker for persisting database entities"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithDependencies[deps](
			info_hash_blocker.Ref,
			logging.Ref,
			plugin_metrics.Ref,
			migrator.Ref,
			plugin_worker.Ref,
			postgres.Ref,
		),
		builder.WithConfig[deps](Ref.MustSub("max_size"), persister.ParamMaxSize),
		builder.WithConfig[deps](Ref.MustSub("max_wait"), persister.ParamMaxWait),
		builder.WithFxOption[deps](
			fx.Provide(
				fx.Annotate(
					func(
						maxSize *atomic.Value[persister.MaxSize],
						maxWait *atomic.Value[persister.MaxWait],
						daoProvider internal_database.DaoTransactionProvider,
						blockerBlocker blocker.Blocker,
						logger *zap.Logger,
						metrics *metrics.Registry,
					) persister.Persister {
						return persister.New(
							maxSize,
							maxWait,
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
		builder.WithWorker(
			func(deps deps) (runner.Provider, worker.Option) {
				return deps.Worker, worker.WithDependencies(
					postgres.Ref,
					info_hash_blocker.Ref,
					migrator.Ref,
				)
			},
		),
	)
)
