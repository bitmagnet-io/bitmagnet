package telemetry

import (
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type DecoratorParams struct {
	fx.In
	DB *gorm.DB
}

type DecoratorResult struct {
	fx.Out
	DB *gorm.DB
}

func NewDecorator(p DecoratorParams) (DecoratorResult, error) {
	db := p.DB
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return DecoratorResult{}, err
	}
	return DecoratorResult{DB: db}, nil
}
