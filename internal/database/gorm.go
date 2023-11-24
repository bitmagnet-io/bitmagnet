package database

import (
	"database/sql"
	gorm2 "github.com/bitmagnet-io/bitmagnet/internal/database/gorm"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

type Params struct {
	fx.In
	Logger    *zap.SugaredLogger
	Lifecycle fx.Lifecycle
	Dialector gorm.Dialector
}

type Result struct {
	fx.Out
	SqlDb  *sql.DB
	GormDb *gorm.DB
}

func New(p Params) (r Result, err error) {
	loggerResult, loggerErr := gorm2.New(gorm2.Params{
		ZapLogger: p.Logger,
		Config: gorm2.Config{
			LogLevel:      gormlogger.Info,
			SlowThreshold: time.Second * 3,
		},
	})
	if loggerErr != nil {
		err = loggerErr
		return
	}
	gDb, dbErr := gorm.Open(p.Dialector, &gorm.Config{
		Logger:               loggerResult.GormLogger,
		DisableAutomaticPing: true,
	})
	if dbErr != nil {
		err = dbErr
		return
	}
	sqlDb, sqlDbErr := gDb.DB()
	if sqlDbErr != nil {
		err = sqlDbErr
		return
	}
	r.GormDb = gDb
	r.SqlDb = sqlDb
	return
}
