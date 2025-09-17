package dev

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/database/migrator"
	"github.com/bitmagnet-io/bitmagnet/internal/env"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/app"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewDevCommand() plugin.Command {
	return func(appBuilder app.Builder) cmd.Command {
		return &DevCommand{
			App: cmd.NewApp[DevCommandDeps](appBuilder),
		}
	}
}

type DevCommandDeps struct {
	fx.In
	DBProvider database.RunnerProvider
	Logger     *zap.Logger
}

type DevCommand struct {
	cmd.Cmd
	cmd.App[DevCommandDeps]
}

func (c *DevCommand) Subcommands() []cmd.Command {
	return []cmd.Command{
		&GormGenCommand{
			App: c.App,
		},
		&MigrateCommand{
			App: c.App,
		},
	}
}

type GormGenCommand struct {
	cmd.Cmd
	cmd.App[DevCommandDeps]
}

func (c *GormGenCommand) Run(env env.Env) error {
	return c.NewRunner(func(deps DevCommandDeps) runner.Runner {
		return runner.SimpleRunner(func(ctx context.Context) error {
			ctx, cancel := context.WithCancelCause(ctx)
			defer cancel(nil)

			shutdown, err := deps.DBProvider.Runner()(ctx, cancel)
			if err != nil {
				return err
			}

			gormDB, err := deps.DBProvider.GormDB()
			if err != nil {
				return err
			}

			generator := gen.BuildGenerator(gormDB)
			generator.Execute()

			return shutdown(ctx)
		})
	})(env)
}

type MigrateCommand struct {
	cmd.Cmd
	cmd.App[DevCommandDeps]
}

func (c *MigrateCommand) Run(env env.Env) error {
	return c.NewRunner(func(deps DevCommandDeps) runner.Runner {
		return runner.SimpleRunner(func(ctx context.Context) error {
			ctx, cancel := context.WithCancelCause(ctx)
			defer cancel(nil)

			shutdown, err := deps.DBProvider.Runner()(ctx, cancel)
			if err != nil {
				return err
			}

			sqlDB, err := deps.DBProvider.SQLDB()
			if err != nil {
				return err
			}

			mgr := migrator.New(sqlDB, deps.Logger)

			err = mgr.Up(ctx)
			if err != nil {
				return err
			}

			return shutdown(ctx)
		})
	})(env)
}
