package level

import "go.uber.org/zap/zapcore"

// revive:disable:line-length-limit

//go:generate go run github.com/abice/go-enum --marshal --names --nocase --nocomments --sql --sqlnullstr --values -t ../../gql/enums.gql.tmpl -f level.go

// Level represents a logging level
/* ENUM(debug, info, warn, error) */
type Level string

// ToZapLevel converts the given string to the appropriate zap level value.
func (level Level) ToZapLevel() zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.WarnLevel
	}
}

func FromZapLevel(l zapcore.Level) Level {
	switch l {
	case zapcore.DebugLevel:
		return LevelDebug
	case zapcore.InfoLevel:
		return LevelInfo
	case zapcore.WarnLevel:
		return LevelWarn
	case zapcore.ErrorLevel:
		return LevelError
	default:
		return LevelWarn
	}
}
