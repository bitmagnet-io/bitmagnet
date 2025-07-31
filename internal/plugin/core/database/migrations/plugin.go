package migrations

import (
	"context"
	"database/sql"

	database_internal "github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/migrations"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	config struct{}

	deps struct {
		fx.In
		database_internal.SQLDBProvider
		Logger *zap.SugaredLogger
	}
)

var (
	Ref = database.Ref.MustSub("migrations")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[config, deps](),
		builder.WithDependencies[config, deps](
			postgres.Ref,
		),
		builder.WithWorkerRegistryOption(
			func(_ config, deps deps) registry.Option {
				return registry.WithWorker(
					Ref.String(),
					deps.runner(),
					worker.WithDependencies(postgres.Ref.String()),
				)
			},
		),
	)
)

func (d deps) runner() runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		var (
			db  *sql.DB
			err error
		)

		defer func() {
			cancel(err)
		}()

		db, err = d.SQLDB()
		if err != nil {
			return nil, err
		}

		migrator := migrations.New(db, d.Logger.Named(Ref.String()))

		err = migrator.Up(ctx)
		if err != nil {
			return nil, err
		}

		return runner.NopShutdowner, nil
	}
}
