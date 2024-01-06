package database

import (
	"database/sql"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
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
	Dialector gorm.Dialector
}

type Result struct {
	fx.Out
	GormDb lazy.Lazy[*gorm.DB]
	SqlDb  lazy.Lazy[*sql.DB]
}

func New(p Params) Result {
	gormDb := lazy.New(func() (*gorm.DB, error) {
		gDb, dbErr := gorm.Open(p.Dialector, &gorm.Config{
			Logger: gorm2.New(gorm2.Params{
				ZapLogger: p.Logger,
				Config: gorm2.Config{
					LogLevel:      gormlogger.Info,
					SlowThreshold: time.Second * 3,
				},
			}).GormLogger,
		})
		if dbErr != nil {
			return nil, dbErr
		}
		return gDb, nil
	})
	return Result{
		GormDb: gormDb,
		SqlDb: lazy.New(func() (*sql.DB, error) {
			gDb, gDbErr := gormDb.Get()
			if gDbErr != nil {
				return nil, gDbErr
			}
			sqlDb, sqlDbErr := gDb.DB()
			if sqlDbErr != nil {
				return nil, sqlDbErr
			}
			return sqlDb, nil
		}),
	}
}
