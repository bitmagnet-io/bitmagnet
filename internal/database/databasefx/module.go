package databasefx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/cache"
	"github.com/bitmagnet-io/bitmagnet/internal/database/healthcheck"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New() fx.Option {
	return fx.Module(
		"database",
		configfx.NewConfigModule[database.Config]("postgres", database.NewDefaultConfig()),
		configfx.NewConfigModule[cache.Config]("gorm_cache", cache.NewDefaultConfig()),
		fx.Provide(
			fx.Annotate(
				healthcheck.New,
				fx.ResultTags(`group:"health_check_options"`),
			),
			fx.Annotate(
				func(
					postgresCfg database.Config,
					cacheCfg cache.Config,
					logger *zap.SugaredLogger,
				) database.RunnerProvider {
					return database.New(
						postgresCfg.CreateDSN(),
						cacheCfg,
						logger,
					)
				},
				fx.As(new(database.PoolProvider)),
				fx.As(new(database.SQLDBProvider)),
				fx.As(new(database.GormDBProvider)),
				fx.As(new(database.DaoProvider)),
				fx.As(new(database.DaoTransactionProvider)),
				fx.As(new(database.Provider)),
				fx.As(fx.Self()),
			),
			fx.Annotate(
				func(prv database.RunnerProvider) registry.Option {
					return registry.WithWorker(
						database.Namespace,
						prv.Runner,
						worker.WithDependencies("logger"),
					)
				},
				fx.ResultTags(`group:"worker_options"`),
			),
			search.New,
		),
	)
}
