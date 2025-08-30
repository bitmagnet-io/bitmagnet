package console

import (
	"github.com/bitmagnet-io/bitmagnet/internal/env"
	"github.com/bitmagnet-io/bitmagnet/internal/logging/encoder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
)

type deps struct {
	fx.In
	Stdout   env.Stdout
	LogLevel zapcore.LevelEnabler
}

var (
	Ref = logging.Ref.MustSub("console")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Outputs logs with console formatting"),
		builder.WithActivation[deps](plugin.ActivationEnabled),
		builder.WithDependencies[deps](logging.Ref),
		builder.WithZapCore(func(deps deps) zapcore.Core {
			return zapcore.NewCore(
				encoder.NewConsole(),
				zapcore.AddSync(deps.Stdout),
				deps.LogLevel,
			)
		}),
	)
)
