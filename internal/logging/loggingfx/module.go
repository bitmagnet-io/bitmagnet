package loggingfx

import (
	"context"
	"errors"

	"github.com/bitmagnet-io/bitmagnet/internal"
	"github.com/bitmagnet-io/bitmagnet/internal/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/logging"
	"github.com/bitmagnet-io/bitmagnet/internal/logging/encoder"
	"github.com/bitmagnet-io/bitmagnet/internal/logging/filerotator"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() fx.Option {
	return fx.Module(
		"logging",
		configfx.NewConfigModule[logging.Config]("log", logging.NewDefaultConfig()),
		fx.Provide(
			func(cfg logging.Config) zapcore.Encoder {
				if cfg.JSON {
					return encoder.NewJSON()
				}

				return encoder.NewConsole()
			},
			func(cfg logging.Config) zap.AtomicLevel {
				return zap.NewAtomicLevelAt(cfg.Level.ToZapLevel())
			},
			func(cfg logging.Config) filerotator.Config {
				return cfg.FileRotator
			},
			func(cfg filerotator.Config) *filerotator.FileRotator {
				return filerotator.New(cfg)
			},
			func(
				stdout internal.Stdout,
				enc zapcore.Encoder,
				lvl zap.AtomicLevel,
				fileRotatorCfg filerotator.Config,
				fileRotator *filerotator.FileRotator,
			) zapcore.Core {
				core := zapcore.NewCore(
					enc,
					zapcore.AddSync(stdout),
					lvl,
				)

				if fileRotatorCfg.Enabled {
					core = zapcore.NewTee(
						core,
						zapcore.NewCore(
							encoder.NewJSON(),
							fileRotator,
							fileRotatorCfg.Level.ToZapLevel(),
						),
					)
				}

				return core
			},
			func(cfg logging.Config, core zapcore.Core) *zap.Logger {
				opts := []zap.Option{
					zap.AddStacktrace(zapcore.ErrorLevel),
					zap.AddCaller(),
				}
				if cfg.Development {
					opts = append(opts, zap.Development())
				}

				return zap.New(core, opts...)
			},
			func(logger *zap.Logger) *zap.SugaredLogger {
				return logger.Sugar()
			},
			fx.Annotate(
				func(
					fileRotatorCfg filerotator.Config,
					fileRotator *filerotator.FileRotator,
				) registry.Option {
					rnr := runner.Nop

					if fileRotatorCfg.Enabled {
						rnr = fileRotator.Runner
					}

					return registry.WithWorker("logger", rnr)
				},
				fx.ResultTags(`group:"worker_options"`),
			),
		),
	)
}

func WithLogger() fx.Option {
	return fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		l := &fxevent.ZapLogger{Logger: log.Named("fx")}
		l.UseLogLevel(zapcore.DebugLevel)

		return fxLogger{l}
	})
}

type fxLogger struct {
	fxevent.Logger
}

func (l fxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.Started:
		if errors.Is(e.Err, context.Canceled) {
			return
		}
	case *fxevent.RollingBack:
		if errors.Is(e.StartErr, context.Canceled) {
			return
		}
	}

	l.Logger.LogEvent(event)
}
