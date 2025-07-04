package migratecmd

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/app"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Command = &cli.Command{
	Name:  "migrate",
	Usage: "Runs database migrations",
	Commands: []*cli.Command{
		{
			Name: "up",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "version",
					Value: 0,
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fxApp := app.New(
					ctx,
					cmd.Writer,
					fx.Invoke(runner(func(migrator migrations.Migrator) error {
						version := cmd.Int64("version")

						if version == 0 {
							return migrator.Up(ctx)
						}

						return migrator.UpTo(ctx, version)
					})),
				)

				fxApp.Run()

				return fxApp.Err()
			},
		},
		{
			Name: "down",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "version",
					Value: 0,
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fxApp := app.New(
					ctx,
					cmd.Writer,
					fx.Invoke(runner(func(migrator migrations.Migrator) error {
						version := cmd.Int64("version")

						if version == 0 {
							return migrator.Down(ctx)
						}

						return migrator.DownTo(ctx, version)
					})),
				)

				fxApp.Run()

				return fxApp.Err()
			},
		},
	},
}

func runner(
	fn func(migrations.Migrator) error,
) func(
	ctx context.Context,
	cfg database.Config,
	logger *zap.SugaredLogger,
) error {
	return func(
		ctx context.Context,
		cfg database.Config,
		logger *zap.SugaredLogger,
	) error {
		pool, err := pgxpool.New(ctx, cfg.CreateDSN())
		if err != nil {
			return err
		}

		defer pool.Close()

		db := stdlib.OpenDBFromPool(pool)

		return fn(migrations.New(db, logger))
	}
}
