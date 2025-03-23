package cli

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/version"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sort"
	"strings"
)

type Params struct {
	fx.In
	Args       []string `name:"cli_args"`
	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Commands   []*cli.Command `group:"commands"`
	Logger     *zap.SugaredLogger
}

type Result struct {
	fx.Out
	App *cli.App
}

func New(p Params) (Result, error) {
	commands := p.Commands
	sort.Slice(commands, func(i, j int) bool {
		return strings.Compare(commands[i].Name, commands[j].Name) < 0
	})
	name := "bitmagnet"
	if version.GitTag != "" {
		name += " " + version.GitTag
	}
	app := &cli.App{
		Name:     name,
		Commands: commands,
		CommandNotFound: func(_ *cli.Context, command string) {
			p.Logger.Errorw("command not found", "command", command)
		},
		After: func(ctx *cli.Context) error {
			return p.Shutdowner.Shutdown()
		},
		Version: version.GitTag,
		// disabling the version flag as for some reason the After hook never gets called
		HideVersion: true,
	}
	app.Setup()
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go (func() {
				// the following hack fixes a weird bug where the CLI does not terminate when calling with just --help
				args := p.Args
				if len(args) == 2 && (args[1] == "-h" || args[1] == "--help") {
					args = []string{args[0]}
				}
				if err := app.RunContext(context.Background(), args); err != nil {
					panic(err)
				}
			})()
			return nil
		},
	})
	return Result{
		App: app,
	}, nil
}
