package logging

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Params struct {
	fx.In
	Config Config
}

type Result struct {
	fx.Out
	Logger      *zap.Logger
	Sugar       *zap.SugaredLogger
	AtomicLevel zap.AtomicLevel
	AppHook     fx.Hook `group:"app_hooks"`
}

func New(params Params) (Result, error) {
	var appHook fx.Hook
	var encoder zapcore.Encoder
	if params.Config.Json {
		encoder = zapcore.NewJSONEncoder(jsonEncoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(consoleEncoderConfig)
	}
	var writeSyncer zapcore.WriteSyncer
	if params.Config.FileRotator.Enabled {
		writeSyncer = NewFileRotator(params.Config.FileRotator)
		appHook = fx.Hook{
			OnStop: func(context.Context) error {
				return writeSyncer.Sync()
			},
		}
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}
	opts := []zap.Option{
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller(),
	}
	if params.Config.Development {
		opts = append(opts, zap.Development())
	}
	atomicLevel := zap.NewAtomicLevelAt(levelToZapLevel(params.Config.Level))
	l := zap.New(zapcore.NewCore(
		encoder,
		writeSyncer,
		atomicLevel,
	), opts...)
	sugar := l.Sugar()
	return Result{
		Logger:      l,
		Sugar:       sugar,
		AtomicLevel: atomicLevel,
		AppHook:     appHook,
	}, nil
}
