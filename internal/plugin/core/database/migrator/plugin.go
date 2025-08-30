package migrator

import (
	"context"
	"database/sql"

	database_internal "github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/migrator"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type deps struct {
	fx.In
	database_internal.SQLDBProvider
	Logger *zap.Logger
}

var (
	Ref = database.Ref.MustSub("migrator")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps](
			"Runs database migrations on startup",
		),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithDependencies[deps](
			postgres.Ref,
		),
		builder.WithWorker(
			func(deps deps) (runner.Provider, worker.Option) {
				return deps.runner(), worker.Options(
					worker.WithDependencies(postgres.Ref),
					worker.ShortLived(),
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

		migrator := migrator.New(db, d.Logger.Named(Ref.String()))

		err = migrator.Up(ctx)
		if err != nil {
			return nil, err
		}

		return runner.NopShutdowner, nil
	}
}
