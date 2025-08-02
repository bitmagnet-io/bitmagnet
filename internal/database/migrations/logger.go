package migrations

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type gooseLogger struct {
	l *zap.Logger
}

func (l gooseLogger) Fatalf(format string, v ...interface{}) {
	l.l.Error(fmt.Sprintf(strings.TrimSpace(format), v...))
}

func (l gooseLogger) Printf(format string, v ...interface{}) {
	fn := l.l.Debug

	if strings.HasPrefix(format, "goose: successfully migrated") ||
		strings.HasPrefix(format, "goose: no migrations to run") {
		fn = l.l.Info
	}

	format = strings.TrimPrefix(format, "goose: ")
	format = strings.TrimSpace(format)

	fn(fmt.Sprintf(format, v...))
}
