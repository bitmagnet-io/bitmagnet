package client

import (
	"github.com/bitmagnet-io/bitmagnet/internal/queue/redis"
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Redis *redis.Client
}

type Result struct {
	fx.Out
	Client *asynq.Client
}

func New(p Params) (Result, error) {
	client := asynq.NewClient(redis.Wrapper{Redis: p.Redis})
	return Result{Client: client}, nil
}
