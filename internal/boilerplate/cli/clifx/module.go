package clifx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/cli"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/cli/args"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"cli",
		fx.Provide(args.New),
		fx.Provide(cli.New),
	)
}
