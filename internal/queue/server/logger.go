package server

import (
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type loggerWrapper struct {
	l *zap.SugaredLogger
}

func NewFromZap(l *zap.SugaredLogger) asynq.Logger {
	return loggerWrapper{l}
}

func (l loggerWrapper) Debug(args ...interface{}) {
	l.l.Debug(args...)
}

func (l loggerWrapper) Info(args ...interface{}) {
	l.l.Info(args...)
}

func (l loggerWrapper) Warn(args ...interface{}) {
	l.l.Warn(args...)
}

func (l loggerWrapper) Error(args ...interface{}) {
	l.l.Error(args...)
}

func (l loggerWrapper) Fatal(args ...interface{}) {
	l.l.Fatal(args...)
}
