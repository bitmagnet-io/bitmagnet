package i18n

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/env"
	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/app"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"go.uber.org/fx"
)

func NewCommand(appBuilder app.Builder) cmd.Command {
	return &i18nCommand{
		App: cmd.NewApp[i18nCommandDeps](appBuilder),
	}
}

type i18nCommand struct {
	cmd.Cmd
	cmd.App[i18nCommandDeps]
}

type i18nCommandDeps struct {
	fx.In
	Provider i18n.Provider
}

func (c *i18nCommand) Name() string {
	return "i18n"
}

func (c *i18nCommand) Subcommands() []cmd.Command {
	return []cmd.Command{
		&extractCommand{
			Cmd: c.Cmd,
			App: c.App,
		},
	}
}

// todo: Should be a dev command
type extractCommand struct {
	cmd.Cmd
	cmd.App[i18nCommandDeps]
}

func (cmd *extractCommand) Run(env env.Env) error {
	return cmd.NewRunner(func(deps i18nCommandDeps) runner.Runner {
		return runner.SimpleRunner(func(context.Context) error {
			return i18n.Write(deps.Provider())
		})
	})(env)
}
