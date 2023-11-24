package migrations

import (
	"go.uber.org/zap"
	"strings"
)

type gooseLogger struct {
	l *zap.SugaredLogger
}

func (l gooseLogger) Fatal(v ...interface{}) {
	l.l.Fatal(v...)
}

func (l gooseLogger) Fatalf(format string, v ...interface{}) {
	l.l.Fatalf(strings.TrimSpace(format), v...)
}

func (l gooseLogger) Print(v ...interface{}) {
	l.l.Debug(v...)
}

func (l gooseLogger) Println(v ...interface{}) {
	l.l.Debug(v...)
}

func (l gooseLogger) Printf(format string, v ...interface{}) {
	fn := l.l.Debugf
	if strings.HasPrefix(format, "goose: successfully migrated") {
		fn = l.l.Infof
	}
	fn(strings.TrimSpace(format), v...)
}
