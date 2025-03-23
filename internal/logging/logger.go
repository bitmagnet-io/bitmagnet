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
	Logger  *zap.Logger
	Sugar   *zap.SugaredLogger
	AppHook fx.Hook `group:"app_hooks"`
}

func New(params Params) Result {
	var appHook fx.Hook
	var encoder zapcore.Encoder
	if params.Config.Json {
		encoder = zapcore.NewJSONEncoder(jsonEncoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(consoleEncoderConfig)
	}
	writeSyncer := zapcore.AddSync(os.Stdout)
	opts := []zap.Option{
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller(),
	}
	if params.Config.Development {
		opts = append(opts, zap.Development())
	}
	core := zapcore.NewCore(
		encoder,
		writeSyncer,
		levelToZapLevel(params.Config.Level),
	)
	if params.Config.FileRotator.Enabled {
		fWriteSyncer := newFileRotator(params.Config.FileRotator)
		core = zapcore.NewTee(
			core,
			zapcore.NewCore(
				zapcore.NewJSONEncoder(jsonEncoderConfig),
				fWriteSyncer,
				levelToZapLevel(params.Config.FileRotator.Level),
			),
		)
		appHook = fx.Hook{
			OnStop: func(context.Context) error {
				return fWriteSyncer.Close()
			},
		}
	}
	l := zap.New(core, opts...)
	return Result{
		Logger:  l,
		Sugar:   l.Sugar(),
		AppHook: appHook,
	}
}
