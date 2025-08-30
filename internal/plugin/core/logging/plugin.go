package logging

import (
	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/logging/level"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type deps struct {
	fx.In
}

var (
	Ref = core.Ref.MustSub("logging")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides logging services"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithDependencies[deps](
			config.Ref,
		),
		builder.WithConfig[deps](Ref.MustSub("level"), level.Param),
		builder.WithFxOption[deps](
			fx.Provide(
				fx.Annotate(
					func(lvl *atomic.Value[level.Level]) zap.AtomicLevel {
						zapLevel := zap.NewAtomicLevel()
						lvl.Subscribe(func(lvl level.Level) {
							zapLevel.SetLevel(lvl.ToZapLevel())
						})
						return zapLevel
					},
					fx.As(new(zapcore.LevelEnabler)),
					fx.As(fx.Self()),
				),
				fx.Annotate(
					func(
						zapCores []zapcore.Core,
					) zapcore.Core {
						return zapcore.NewTee(zapCores...)
					},
					fx.ParamTags(`group:"zap_cores"`),
				),
				func(core zapcore.Core) *zap.Logger {
					opts := []zap.Option{
						zap.AddStacktrace(zapcore.ErrorLevel),
						zap.AddCaller(),
					}

					return zap.New(core, opts...)
				},
			),
			// fx.Invoke(func(
			// 	logger *zap.Logger,
			// 	availablePlugins plugin.Available,
			// 	enabledPlugins plugin.Enabled,
			// ) {
			// 	logger = logger.Named("plugins")
			// 	logger.Info(
			// 		"started",
			// 		zap.Strings("enabled", ref.Refs(enabledPlugins).Strings()),
			// 		zap.Strings("available", ref.Refs(availablePlugins).Strings()),
			// 	)
			// }),
		),
	)
)
