package console

import (
	"github.com/bitmagnet-io/bitmagnet/internal/env"
	"github.com/bitmagnet-io/bitmagnet/internal/logging/encoder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
)

type (
	config struct{}

	deps struct {
		fx.In
		Stdout   env.Stdout
		LogLevel zapcore.LevelEnabler
	}
)

var (
	Ref = logging.Ref.MustSub("console")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[config, deps](),
		builder.WithDependencies[config, deps](logging.Ref),
		builder.WithZapCore(func(config config, deps deps) zapcore.Core {
			return zapcore.NewCore(
				encoder.NewConsole(),
				zapcore.AddSync(deps.Stdout),
				deps.LogLevel,
			)
		}),
	)
)
