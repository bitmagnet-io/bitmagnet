package database

import (
	"database/sql"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

type Params struct {
	fx.In
	SqlDB  lazy.Lazy[*sql.DB]
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	GormDb lazy.Lazy[*gorm.DB]
}

func New(p Params) Result {
	return Result{
		GormDb: lazy.New(func() (*gorm.DB, error) {
			sqlDB, err := p.SqlDB.Get()
			if err != nil {
				return nil, err
			}
			dialector := postgres.New(postgres.Config{
				Conn: sqlDB,
			})
			gDb, dbErr := gorm.Open(dialector, &gorm.Config{
				Logger: logger.New(logger.Params{
					ZapLogger: p.Logger,
					Config: logger.Config{
						LogLevel:      gormlogger.Info,
						SlowThreshold: time.Second * 3,
					},
				}).GormLogger,
			})
			if dbErr != nil {
				return nil, dbErr
			}
			return gDb, nil
		}),
	}
}
