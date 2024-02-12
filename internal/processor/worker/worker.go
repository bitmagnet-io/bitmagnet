package worker

import (
	"context"
	"encoding/json"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/worker"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

type Params struct {
	fx.In
	Neoq      lazy.Lazy[queue.Queue]
	Processor lazy.Lazy[processor.Processor]
	Logger    *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Worker worker.Worker `group:"workers"`
}

func New(p Params) (Result, error) {
	var q queue.Queue
	return Result{
		Worker: worker.NewWorker(
			"queue_server",
			fx.Hook{
				OnStart: func(ctx context.Context) error {
					q, err := p.Neoq.Get()
					if err != nil {
						return err
					}
					pr, err := p.Processor.Get()
					if err != nil {
						return err
					}
					return q.Start(ctx, handler.New(processor.MessageName, func(ctx context.Context, job model.QueueJob) (err error) {
						msg := &processor.MessageParams{}
						if err := json.Unmarshal([]byte(job.Payload), msg); err != nil {
							return err
						}
						return pr.Process(ctx, *msg)
					}, handler.JobTimeout(time.Second*60*5), handler.Concurrency(3)))
				},
				OnStop: func(ctx context.Context) error {
					if q != nil {
						q.Shutdown(ctx)
					}
					return nil
				},
			},
		),
	}, nil
}
