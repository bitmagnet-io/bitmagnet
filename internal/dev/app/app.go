package app

import (
	"github.com/bitmagnet-io/bitmagnet/internal/dev/devfx"
	"github.com/bitmagnet-io/bitmagnet/internal/logging/loggingfx"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New() *fx.App {
	return fx.New(
		devfx.New(),
		loggingfx.WithLogger(),
		fx.Invoke(func(
			_ *cli.App,
			logger *zap.SugaredLogger,
		) {
			logger.Debug("dev invoked")
		}),
	)
}
