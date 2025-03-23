package migrations

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DecoratorParams struct {
	fx.In
	DB       lazy.Lazy[*gorm.DB]
	Migrator lazy.Lazy[Migrator]
	Logger   *zap.SugaredLogger
}

type DecoratorResult struct {
	fx.Out
	DB lazy.Lazy[*gorm.DB]
}

func NewDecorator(p DecoratorParams) DecoratorResult {
	return DecoratorResult{
		DB: lazy.New(func() (*gorm.DB, error) {
			db, err := p.DB.Get()
			if err != nil {
				return nil, err
			}
			m, err := p.Migrator.Get()
			if err != nil {
				return nil, err
			}
			if migrateErr := m.Up(context.TODO()); migrateErr != nil {
				return nil, migrateErr
			}
			return db, nil
		}),
	}
}
