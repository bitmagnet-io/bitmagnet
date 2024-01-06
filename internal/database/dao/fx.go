package dao

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Params struct {
	fx.In
	GormDb lazy.Lazy[*gorm.DB]
}

type Result struct {
	fx.Out
	Dao lazy.Lazy[*Query]
}

func New(p Params) Result {
	return Result{
		Dao: lazy.New(func() (*Query, error) {
			db, err := p.GormDb.Get()
			if err != nil {
				return nil, err
			}
			return Use(db), nil
		}),
	}
}
