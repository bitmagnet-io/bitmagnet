package queuemetrics

import (
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Params struct {
	fx.In
	DB lazy.Lazy[*gorm.DB]
}

type Result struct {
	fx.Out
	Client lazy.Lazy[Client]
}

func New(p Params) Result {
	return Result{
		Client: lazy.New[Client](func() (Client, error) {
			db, err := p.DB.Get()
			if err != nil {
				return nil, err
			}
			return client{db}, nil
		}),
	}
}
