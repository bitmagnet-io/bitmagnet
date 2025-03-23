package database

import (
	"database/sql"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database/exclause"
	"github.com/bitmagnet-io/bitmagnet/internal/database/logger"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Params struct {
	fx.In
	SQLDB  lazy.Lazy[*sql.DB]
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	GormDB lazy.Lazy[*gorm.DB]
}

func New(p Params) Result {
	return Result{
		GormDB: lazy.New(func() (*gorm.DB, error) {
			sqlDB, err := p.SQLDB.Get()
			if err != nil {
				return nil, err
			}
			dialector := postgres.New(postgres.Config{
				Conn: sqlDB,
			})
			gDB, dbErr := gorm.Open(dialector, &gorm.Config{
				Logger: logger.New(logger.Params{
					ZapLogger: p.Logger,
					Config: logger.Config{
						LogLevel:      gormlogger.Info,
						SlowThreshold: time.Second * 30,
					},
				}).GormLogger,
				DisableAutomaticPing: true,
			})
			if dbErr != nil {
				return nil, dbErr
			}
			if pluginErr := gDB.Use(exclause.New()); pluginErr != nil {
				return nil, pluginErr
			}
			return gDB, nil
		}),
	}
}
