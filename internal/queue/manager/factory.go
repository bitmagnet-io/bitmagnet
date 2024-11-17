package manager

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Params struct {
	fx.In
	Dao lazy.Lazy[*dao.Query]
	DB  lazy.Lazy[*gorm.DB]
}

type Result struct {
	fx.Out
	Manager lazy.Lazy[Manager]
}

func New(params Params) Result {
	return Result{
		Manager: lazy.New[Manager](func() (Manager, error) {
			d, err := params.Dao.Get()
			if err != nil {
				return nil, err
			}
			db, err := params.DB.Get()
			if err != nil {
				return nil, err
			}
			return manager{
				dao: d,
				db:  db,
			}, nil
		}),
	}
}
