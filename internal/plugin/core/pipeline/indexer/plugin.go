package indexer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/indexer"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"go.uber.org/fx"
)

type (
	config struct{}

	deps struct {
		fx.In
		Indexer indexer.Indexer
	}
)

var (
	Ref = pipeline.Ref.MustSub("processor")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithFxOption[config, deps](
			fx.Provide(
				indexer.New,
				func() queue.JobProvider[indexer.MessageParams] {
					return func(msg indexer.MessageParams, options ...model.QueueJobOption) (model.QueueJob, error) {
						return model.NewQueueJob(
							Ref.String(),
							msg,
							append([]model.QueueJobOption{model.QueueJobMaxRetries(2)}, options...)...,
						)
					}
				},
			),
		),
		builder.WithQueueHandler(
			func(cfg config, deps deps) handler.Handler {
				return handler.New(
					Ref.String(),
					func(job model.QueueJob) runner.Runner {
						return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
							msg := &indexer.MessageParams{}
							if err := json.Unmarshal([]byte(job.Payload), msg); err != nil {
								return runner.NopShutdowner, err
							}

							return deps.Indexer.NewJob(*msg)(ctx, cancel)
						}
					},
					handler.JobTimeout(time.Second*60*10),
					handler.Concurrency(1),
				)
			},
		),
	)
)
