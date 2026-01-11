package json

import (
	"github.com/bitmagnet-io/bitmagnet/internal/logging/encoder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
)

type deps struct {
	fx.In
	Writer env.Stdout
	Level  zapcore.LevelEnabler
}

var (
	Ref = logging.Ref.MustSub("json")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Outputs logs in JSON format"),
		builder.WithActivation[deps](plugin.ActivationDisabled),
		builder.WithDependencies[deps](logging.Ref),
		builder.WithZapCore(func(deps deps) zapcore.Core {
			return zapcore.NewCore(
				encoder.NewJSON(),
				zapcore.AddSync(deps.Writer),
				deps.Level,
			)
		}),
	)
)
