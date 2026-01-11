package dev

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/database/migrator"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
)

type MigrateCommand struct {
	cmd.Cmd `cmd:"doc=Run database migrations"`
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
