package migrations

import (
	"context"
	"database/sql"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type DecoratorParams struct {
	fx.In
	DB     lazy.Lazy[*sql.DB]
	Logger *zap.SugaredLogger
}

type DecoratorResult struct {
	fx.Out
	DB lazy.Lazy[*sql.DB]
}

func NewDecorator(p DecoratorParams) DecoratorResult {
	return DecoratorResult{
		DB: lazy.New(func() (*sql.DB, error) {
			db, err := p.DB.Get()
			if err != nil {
				return nil, err
			}
			m := New(Params{
				DB:     db,
				Logger: p.Logger,
			})
			if migrateErr := m.Up(context.TODO()); migrateErr != nil {
				return nil, migrateErr
			}
			return db, nil
		}),
	}
}
