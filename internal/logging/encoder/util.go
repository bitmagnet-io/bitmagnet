package encoder

import (
	"strings"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/colors"
	"github.com/bitmagnet-io/bitmagnet/internal/logging/level"
	"github.com/fatih/color"
	"go.uber.org/zap/zapcore"
)

const (
	timestamp  = "timestamp"
	severity   = "severity"
	logger     = "logger"
	caller     = "caller"
	message    = "message"
	stacktrace = "stacktrace"
)

// levelEncoder transforms a zap level to the associated stackdriver level.
func levelEncoder() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(level.FromZapLevel(l).String())
	}
}

// timeEncoder encodes the time as RFC3339 nano.
func timeEncoder() zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339Nano))
	}
}

// timeEncoder encodes the time as RFC3339 nano.
func consoleTimeEncoder() zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(colors.Grey.Sprintf("%s", t.Format(time.RFC3339)))
	}
}

func paddedNameEncoder(minLength int) zapcore.NameEncoder {
	return func(name string, enc zapcore.PrimitiveArrayEncoder) {
		if padN := minLength - len(name); padN > 0 {
			name += strings.Repeat(" ", padN)
		}

		enc.AppendString(name)
	}
}

var levelColors = map[level.Level]*color.Color{
	level.LevelDebug: colors.Magenta,
	level.LevelInfo:  colors.Blue,
	level.LevelWarn:  colors.Yellow,
	level.LevelError: colors.Red,
}

func colorLevelEncoder(zapLevel zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	lvl := level.FromZapLevel(zapLevel)
	str := strings.ToUpper(lvl.String())
	str += strings.Repeat(" ", 5-len(str))
	enc.AppendString(levelColors[lvl].Sprintf("%s", str))
}
