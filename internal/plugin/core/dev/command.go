package dev

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/app"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
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
	DBProvider   database.RunnerProvider
	I18nProvider i18n.Provider
	Logger       *zap.Logger
}

type DevCommand struct {
	cmd.Cmd `cmd:"doc=Development utilities"`
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
		&I18nCommand{
			App: c.App,
		},
	}
}
