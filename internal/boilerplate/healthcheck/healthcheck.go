package healthcheck

import (
	"github.com/hellofresh/health-go/v5"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Options []health.Option `group:"healthcheck_options"`
}

type Result struct {
	fx.Out
	Health *health.Health
}

func New(p Params) (Result, error) {
	h, err := health.New(append(p.Options, health.WithSystemInfo())...)
	if err != nil {
		return Result{}, err
	}
	return Result{Health: h}, nil
}
