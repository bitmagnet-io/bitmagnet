package cache

import (
	caches "github.com/mgdigital/gorm-cache/v2"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type DecoratorParams struct {
	fx.In
	Plugin *caches.Caches
	DB     *gorm.DB
}

type DecoratorResult struct {
	fx.Out
	DB *gorm.DB
}

func NewDecorator(p DecoratorParams) (DecoratorResult, error) {
	db := p.DB
	if err := db.Use(p.Plugin); err != nil {
		return DecoratorResult{}, err
	}
	return DecoratorResult{DB: db}, nil
}
