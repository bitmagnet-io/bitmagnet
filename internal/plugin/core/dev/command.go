package dev

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/app"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd"
	"github.com/bitmagnet-io/bitmagnet/pkg/i18n"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewCommand() plugin.Command {
	return func(appBuilder app.Builder) cmd.Command {
		return &Command{
			App: cmd.NewApp[CommandDeps](appBuilder),
		}
	}
}

type CommandDeps struct {
	fx.In
	DBProvider   database.RunnerProvider
	I18nProvider i18n.MessageProvider
	Logger       *zap.Logger
}

type Command struct {
	cmd.Cmd `cmd:"name=dev,doc=Development utilities"`
	cmd.App[CommandDeps]
}

func (c *Command) Subcommands() []cmd.Command {
	return []cmd.Command{
		&GormGenCommand{
			App: c.App,
		},
		&MigrateCommand{
			App: c.App,
		},
		&I18nCommand{
			App: c.App,
		},
	}
}
