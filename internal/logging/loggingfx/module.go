package loggingfx

import (
	"context"
	"errors"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
