package batcher

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/indexer"
	"github.com/bitmagnet-io/bitmagnet/internal/indexer/batch"
	batch_queue "github.com/bitmagnet-io/bitmagnet/internal/indexer/batch/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline"
	plugin_indexer "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline/indexer"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	DaoProvider        database.DaoProvider
	ProcessJobProvider queue.JobProvider[indexer.MessageParams]
	BatchJobProvider   queue.JobProvider[batch.MessageParams]
}

var (
	Ref = pipeline.Ref.MustSub("batcher")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithDependencies[deps](
			postgres.Ref,
			plugin_indexer.Ref,
		),
		builder.WithFxOption[deps](
			fx.Provide(func() queue.JobProvider[batch.MessageParams] {
				return func(msg batch.MessageParams, options ...model.QueueJobOption) (model.QueueJob, error) {
					if msg.BatchSize == 0 {
						msg.BatchSize = 100
					}

					if msg.ChunkSize == 0 {
						msg.ChunkSize = 10_000
					}

					return model.NewQueueJob(
						Ref.String(),
						msg,
						append([]model.QueueJobOption{model.QueueJobMaxRetries(2)}, options...)...,
					)
				}
			}),
		),
		builder.WithQueueHandler(
			func(deps deps) handler.Handler {
				return handler.New(
					Ref.String(),
					batch_queue.New(deps.DaoProvider, deps.ProcessJobProvider, deps.BatchJobProvider),
					handler.JobTimeout(time.Second*60*10),
					handler.Concurrency(1),
				)
			},
		),
	)
)
