package dev

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
)

type I18nCommand struct {
	cmd.Cmd `cmd:"doc=Internationalization commands"`
	cmd.App[CommandDeps]
}

func (*I18nCommand) Name() string {
	return "i18n"
}

func (c *I18nCommand) Subcommands() []cmd.Command {
	return []cmd.Command{
		&extractCommand{
			Cmd: c.Cmd,
			App: c.App,
		},
	}
}

type extractCommand struct {
	cmd.Cmd `cmd:"doc=Extract translatable strings to i18n files"`
	cmd.App[CommandDeps]
}

func (cmd *extractCommand) Run(env env.Env) error {
	return cmd.NewRunner(func(deps CommandDeps) runner.Runner {
		return runner.SimpleRunner(func(context.Context) error {
			return i18n.Write(deps.I18nProvider())
		})
	})(env)
}
