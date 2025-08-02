package logging

import (
	"github.com/bitmagnet-io/bitmagnet/internal/logging/level"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Config struct {
		Level       level.Level
		Development bool
	}

	deps struct{}
)

func NewDefaultConfig() Config {
	return Config{
		Level:       "info",
		Development: false,
	}
}

var (
	Ref = core.Ref.MustSub("logging")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[Config, deps](),
		builder.WithDependencies[Config, deps](
			config.Ref,
		),
		builder.WithDefaultConfig[Config, deps](NewDefaultConfig()),
		builder.WithFxOption[Config, deps](
			fx.Provide(
				fx.Annotate(
					func(cfg Config) zap.AtomicLevel {
						return zap.NewAtomicLevelAt(cfg.Level.ToZapLevel())
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
				func(cfg Config, core zapcore.Core) *zap.Logger {
					opts := []zap.Option{
						zap.AddStacktrace(zapcore.ErrorLevel),
						zap.AddCaller(),
					}
					if cfg.Development {
						opts = append(opts, zap.Development())
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
