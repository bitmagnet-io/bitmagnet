package postgres

import (
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Params struct {
	fx.In
	Config Config
}

type Result struct {
	fx.Out
	Dialector gorm.Dialector
}

func New(p Params) (Result, error) {
	return Result{
		Dialector: postgres.Open(p.Config.DSN()),
	}, nil
}
