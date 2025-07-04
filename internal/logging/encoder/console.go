package encoder

import "go.uber.org/zap/zapcore"

func NewConsole() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(consoleEncoderConfig)
}

var consoleEncoderConfig = zapcore.EncoderConfig{
	TimeKey:          timestamp,
	LevelKey:         "L",
	NameKey:          "N",
	CallerKey:        zapcore.OmitKey,
	FunctionKey:      zapcore.OmitKey,
	MessageKey:       "M",
	StacktraceKey:    "S",
	LineEnding:       zapcore.DefaultLineEnding,
	EncodeLevel:      colorLevelEncoder,
	EncodeTime:       consoleTimeEncoder(),
	EncodeDuration:   zapcore.StringDurationEncoder,
	EncodeCaller:     zapcore.ShortCallerEncoder,
	ConsoleSeparator: "  ",
	EncodeName:       paddedNameEncoder(22),
}
