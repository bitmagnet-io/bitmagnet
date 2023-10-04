package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

type Config struct {
	Level       string
	Development bool
	Json        bool
}

func NewDefault() Config {
	return Config{
		Level:       "info",
		Development: false,
		Json:        false,
	}
}

func NewZapConfig(c Config) zap.Config {
	var zc zap.Config
	if c.Development {
		zc = NewDevelopmentConfig()
	} else {
		zc = NewProductionConfig()
	}
	if c.Level != "" {
		zc.Level = zap.NewAtomicLevelAt(levelToZapLevel(c.Level))
	}
	if c.Json {
		zc.Encoding = encodingJSON
		zc.EncoderConfig = jsonEncoderConfig
	} else {
		zc.Encoding = encodingConsole
		zc.EncoderConfig = consoleEncoderConfig
	}
	return zc
}

func NewDevelopmentConfig() zap.Config {
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         encodingConsole,
		EncoderConfig:    consoleEncoderConfig,
		OutputPaths:      outputStderr,
		ErrorOutputPaths: outputStderr,
	}
}

func NewProductionConfig() zap.Config {
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         encodingJSON,
		EncoderConfig:    jsonEncoderConfig,
		OutputPaths:      outputStderr,
		ErrorOutputPaths: outputStderr,
	}
}

const (
	timestamp  = "timestamp"
	severity   = "severity"
	logger     = "logger"
	caller     = "caller"
	message    = "message"
	stacktrace = "stacktrace"

	levelDebug     = "DEBUG"
	levelInfo      = "INFO"
	levelWarning   = "WARNING"
	levelError     = "ERROR"
	levelCritical  = "CRITICAL"
	levelAlert     = "ALERT"
	levelEmergency = "EMERGENCY"

	encodingConsole = "console"
	encodingJSON    = "json"
)

var jsonEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        timestamp,
	LevelKey:       severity,
	NameKey:        logger,
	CallerKey:      caller,
	MessageKey:     message,
	StacktraceKey:  stacktrace,
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    levelEncoder(),
	EncodeTime:     timeEncoder(),
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

var consoleEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "",
	LevelKey:       "L",
	NameKey:        "N",
	CallerKey:      "C",
	FunctionKey:    zapcore.OmitKey,
	MessageKey:     "M",
	StacktraceKey:  "S",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

// levelToZapLevel converts the given string to the appropriate zap level
// value.
func levelToZapLevel(s string) zapcore.Level {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case levelDebug:
		return zapcore.DebugLevel
	case levelInfo:
		return zapcore.InfoLevel
	case levelWarning:
		return zapcore.WarnLevel
	case levelError:
		return zapcore.ErrorLevel
	case levelCritical:
		return zapcore.DPanicLevel
	case levelAlert:
		return zapcore.PanicLevel
	case levelEmergency:
		return zapcore.FatalLevel
	}

	return zapcore.WarnLevel
}

// levelEncoder transforms a zap level to the associated stackdriver level.
func levelEncoder() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString(levelDebug)
		case zapcore.InfoLevel:
			enc.AppendString(levelInfo)
		case zapcore.WarnLevel:
			enc.AppendString(levelWarning)
		case zapcore.ErrorLevel:
			enc.AppendString(levelError)
		case zapcore.DPanicLevel:
			enc.AppendString(levelCritical)
		case zapcore.PanicLevel:
			enc.AppendString(levelAlert)
		case zapcore.FatalLevel:
			enc.AppendString(levelEmergency)
		}
	}
}

// timeEncoder encodes the time as RFC3339 nano.
func timeEncoder() zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339Nano))
	}
}

var outputStderr = []string{"stderr"}
