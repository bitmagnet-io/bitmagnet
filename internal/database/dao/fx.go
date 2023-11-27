package dao

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Params struct {
	fx.In
	GormDb *gorm.DB
}

type Result struct {
	fx.Out
	Dao *Query
}

func New(p Params) (r Result, err error) {
	r.Dao = Use(p.GormDb)
	//SetDefault(p.GormDb)
	return
}
