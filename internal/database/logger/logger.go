package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Config struct {
	SlowThreshold time.Duration
	LogLevel      gormlogger.LogLevel
}

func New(logger *zap.Logger, cfg Config) gormlogger.Interface {
	return &customLogger{
		logLevel:      cfg.LogLevel,
		slowThreshold: cfg.SlowThreshold,
		zap:           logger,
	}
}

type customLogger struct {
	logLevel      gormlogger.LogLevel
	slowThreshold time.Duration
	zap           *zap.Logger
}

func (l *customLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.logLevel = level

	return l
}

func (l *customLogger) Info(_ context.Context, msg string, data ...interface{}) {
	l.zap.Debug(msg, zap.Any("data", data))
}

func (l *customLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	l.zap.Warn(msg, zap.Any("data", data))
}

func (l *customLogger) Error(_ context.Context, msg string, data ...interface{}) {
	l.zap.Error(msg, zap.Any("data", data))
}

func (l *customLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)

	switch {
	case err != nil && l.logLevel >= gormlogger.Error &&
		!errors.Is(err, gormlogger.ErrRecordNotFound) &&
		!errors.Is(err, context.Canceled):
		sql, rows := fc()
		l.zap.Error("gorm trace",
			zap.String("location", utils.FileWithLineNum()),
			zap.Error(err),
			zap.Float64("elapsed", float64(elapsed.Nanoseconds())/1e6),
			zap.String("sql", sql),
			zap.Int64("rows", rows))
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= gormlogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.slowThreshold)
		l.zap.Warn("gorm trace",
			zap.String("location", utils.FileWithLineNum()),
			zap.String("slowLog", slowLog),
			zap.Float64("elapsed", float64(elapsed.Nanoseconds())/1e6),
			zap.String("sql", sql),
			zap.Int64("rows", rows))
	case l.logLevel == gormlogger.Info:
		sql, rows := fc()
		l.zap.Debug("gorm trace",
			zap.String("location", utils.FileWithLineNum()),
			zap.Float64("elapsed", float64(elapsed.Nanoseconds())/1e6),
			zap.String("sql", sql),
			zap.Int64("rows", rows))
	}
}
