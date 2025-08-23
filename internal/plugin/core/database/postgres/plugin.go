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
		builder.WithEnabledByDefault[deps](),
		builder.WithConfigParam[deps](Ref.MustSub("dsn"), ParamDSN),
		builder.WithConfigParam[deps](Ref.MustSub("host"), ParamHost),
		builder.WithConfigParam[deps](Ref.MustSub("port"), ParamPort),
		builder.WithConfigParam[deps](Ref.MustSub("user"), ParamUser),
		builder.WithConfigParam[deps](Ref.MustSub("password"), ParamPassword),
		builder.WithConfigParam[deps](Ref.MustSub("database"), ParamDatabase),
		builder.WithFxOption[deps](
			fx.Provide(
				fx.Annotate(
					func(
						plugins []gorm.Plugin,
					) []gorm.Plugin {
						return plugins
					},
					fx.ParamTags(`group:"gorm_plugins"`),
				),
				fx.Annotate(
					func(
						plugins []gorm.Plugin,
						cfg Config,
						logger *zap.Logger,
					) database.RunnerProvider {
						return database.New(
							string(cfg.CreateDSN()),
							plugins,
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
			),
		),
		builder.WithWorkerRegistryOption(func(deps deps) registry.Option {
			return registry.WithWorker(
				Ref.String(),
				deps.Provider,
			)
		}),
		builder.WithHealthCheckerOption(func(deps deps) health.CheckerOption {
			return healthcheck.New(Ref.String(), deps.Provider)
		}),
	)
)
