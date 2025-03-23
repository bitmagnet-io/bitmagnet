package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Config struct {
	SlowThreshold time.Duration
	LogLevel      gormlogger.LogLevel
}

type Params struct {
	fx.In
	Config    Config
	ZapLogger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	GormLogger gormlogger.Interface
}

func New(p Params) Result {
	return Result{
		GormLogger: &customLogger{
			logLevel:      p.Config.LogLevel,
			slowThreshold: p.Config.SlowThreshold,
			zap:           p.ZapLogger.Named("gorm"),
		},
	}
}

type customLogger struct {
	logLevel      gormlogger.LogLevel
	slowThreshold time.Duration
	zap           *zap.SugaredLogger
}

func (l *customLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.logLevel = level

	return l
}

func (l *customLogger) Info(_ context.Context, msg string, data ...interface{}) {
	l.zap.Debugw("gorm", "msg", msg, "data", data)
}

func (l *customLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	l.zap.Warnw("gorm", "msg", msg, "data", data)
}

func (l *customLogger) Error(_ context.Context, msg string, data ...interface{}) {
	l.zap.Errorw("gorm", "msg", msg, "data", data)
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
		l.zap.Errorw("gorm trace",
			"location", utils.FileWithLineNum(),
			"error", err,
			"elapsed", float64(elapsed.Nanoseconds())/1e6,
			"sql", sql,
			"rows", rows)
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= gormlogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.slowThreshold)
		l.zap.Warnw("gorm trace",
			"location", utils.FileWithLineNum(),
			"slowLog", slowLog,
			"elapsed", float64(elapsed.Nanoseconds())/1e6,
			"sql", sql,
			"rows", rows)
	case l.logLevel == gormlogger.Info:
		sql, rows := fc()
		l.zap.Debugw("gorm trace",
			"location", utils.FileWithLineNum(),
			"elapsed", float64(elapsed.Nanoseconds())/1e6,
			"sql", sql,
			"rows", rows)
	}
}
