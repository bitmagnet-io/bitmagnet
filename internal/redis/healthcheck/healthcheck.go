package healthcheck

import (
	"context"
	"github.com/hellofresh/health-go/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Redis *redis.Client
}

type Result struct {
	fx.Out
	Option health.Option `group:"healthcheck_options"`
}

func New(p Params) (r Result, err error) {
	r.Option = health.WithChecks(health.Config{
		Name: "redis",
		Check: func(ctx context.Context) error {
			_, err := p.Redis.Ping(ctx).Result()
			return err
		},
	})
	return
}
