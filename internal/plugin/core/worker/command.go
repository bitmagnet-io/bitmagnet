package worker

import (
	"github.com/bitmagnet-io/bitmagnet/internal/env"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/app"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"go.uber.org/fx"
)

func NewStartCommand() plugin.Command {
	return func(appBuilder app.Builder) cmd.Command {
		return &StartCommand{
			App: cmd.NewApp[StartDeps](appBuilder),
		}
	}
}

type StartCommand struct {
	cmd.Cmd
	cmd.App[StartDeps]
}

type StartDeps struct {
	fx.In
	Registry *registry.Registry
}

func (cmd *StartCommand) Run(env env.Env) error {
	return cmd.NewRunner(func(deps StartDeps) runner.Runner {
		return deps.Registry.Runner()
	})(env)
}
