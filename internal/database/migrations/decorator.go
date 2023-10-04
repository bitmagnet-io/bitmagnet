package migrations

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DecoratorParams struct {
	fx.In
	DB     *gorm.DB
	Logger *zap.SugaredLogger
}

type DecoratorResult struct {
	fx.Out
	DB *gorm.DB
}

func NewDecorator(p DecoratorParams) (result DecoratorResult, err error) {
	result.DB = p.DB
	sqlDb, dbErr := p.DB.DB()
	if dbErr != nil {
		err = dbErr
		return
	}
	// avoid failing here on a non-connectable database
	pingErr := sqlDb.Ping()
	if pingErr != nil {
		p.Logger.Errorf("failed to ping database: %v", pingErr)
		return
	}
	m := New(Params{
		DB:     sqlDb,
		Logger: p.Logger,
	})
	migrateErr := m.Up(context.TODO())
	if migrateErr != nil {
		err = migrateErr
		return
	}
	return
}
