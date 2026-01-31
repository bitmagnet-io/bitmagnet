package postgres

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/healthcheck"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type deps struct {
	fx.In
	Provider database.RunnerProvider
}

var (
	Ref = ref.Root.MustSub("postgres")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Initiates and manages the Postgres database connection"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithConfig[deps](Ref.MustSub("dsn"), ParamDSN),
		builder.WithConfig[deps](Ref.MustSub("host"), ParamHost),
		builder.WithConfig[deps](Ref.MustSub("port"), ParamPort),
		builder.WithConfig[deps](Ref.MustSub("user"), ParamUser),
		builder.WithConfig[deps](Ref.MustSub("password"), ParamPassword),
		builder.WithConfig[deps](Ref.MustSub("database"), ParamDatabase),
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
		builder.WithWorker(
			func(deps deps) (runner.Provider, worker.Option) {
				return deps.Provider, nil
			},
		),
		builder.WithHealthCheckerOption(func(deps deps) health.CheckerOption {
			return healthcheck.New(Ref.String(), deps.Provider)
		}),
	)
)
