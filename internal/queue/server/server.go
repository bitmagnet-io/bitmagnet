package server

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/worker"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/consumer"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/redis"
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Config    queue.Config
	Redis     lazy.Lazy[*redis.Client]
	Consumers []lazy.Lazy[consumer.Consumer] `group:"queue_consumers"`
	Options   []Option                       `group:"queue_server_options"`
	Logger    *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Worker worker.Worker `group:"workers"`
}

func New(p Params) (Result, error) {
	var srv *asynq.Server
	return Result{
		Worker: worker.NewWorker(
			"queue_server",
			fx.Hook{
				OnStart: func(ctx context.Context) error {
					cfg := &asynq.Config{
						Concurrency: p.Config.Concurrency,
						Logger:      loggerWrapper{p.Logger.Named("asynq")},
						LogLevel:    asynq.DebugLevel,
						Queues:      p.Config.Queues,
					}
					for _, opt := range p.Options {
						opt.apply(cfg)
					}
					srv = asynq.NewServer(redis.Wrapper{Redis: p.Redis}, *cfg)
					mux := asynq.NewServeMux()
					for _, lc := range p.Consumers {
						c, err := lc.Get()
						if err != nil {
							return err
						}
						mux.Handle(c.Pattern(), c)
					}
					return srv.Start(mux)
				},
				OnStop: func(ctx context.Context) error {
					if srv != nil {
						srv.Shutdown()
					}
					return nil
				},
			},
		),
	}, nil
}

type Option struct {
	apply func(cfg *asynq.Config)
}
