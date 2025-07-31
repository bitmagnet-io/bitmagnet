package postgres

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/healthcheck"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	plugin_database "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type deps struct {
	fx.In
	Provider database.RunnerProvider
}

var (
	Ref = plugin_database.Ref.MustSub("postgres")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[Config, deps](),
		builder.WithDefaultConfig[Config, deps](NewDefaultConfig()),
		builder.WithFxOption[Config, deps](
			fx.Provide(
				fx.Annotate(
					func(
						plugins []gorm.Plugin,
						cfg Config,
						logger *zap.Logger,
					) database.RunnerProvider {
						return database.New(
							cfg.CreateDSN(),
							plugins,
							logger,
						)
					},
					fx.ParamTags(`group:"gorm_plugins"`),
					fx.As(new(database.PoolProvider)),
					fx.As(new(database.SQLDBProvider)),
					fx.As(new(database.GormDBProvider)),
					fx.As(new(database.DaoProvider)),
					fx.As(new(database.DaoTransactionProvider)),
					fx.As(new(database.Provider)),
					fx.As(fx.Self()),
				),
			),
		),
		builder.WithWorkerRegistryOption(func(_ Config, deps deps) registry.Option {
			return registry.WithWorker(
				Ref.String(),
				deps.Provider,
			)
		}),
		builder.WithHealthCheckerOption(func(_ Config, deps deps) health.CheckerOption {
			return healthcheck.New(Ref.String(), deps.Provider)
		}),
	)
)
