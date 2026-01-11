package dev

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/database/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
)

type GormGenCommand struct {
	cmd.Cmd `cmd:"doc=Generate GORM models from the database schema"`
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
