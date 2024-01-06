package client

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/redis"
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Redis lazy.Lazy[*redis.Client]
}

type Result struct {
	fx.Out
	Client lazy.Lazy[*asynq.Client]
}

func New(p Params) Result {
	return Result{
		Client: lazy.New(func() (*asynq.Client, error) {
			return asynq.NewClient(redis.Wrapper{Redis: p.Redis}), nil
		}),
	}
}
