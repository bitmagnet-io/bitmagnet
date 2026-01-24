package batcher

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/postgres"
	plugin_indexer "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/batch"
	batch_queue "github.com/bitmagnet-io/bitmagnet/internal/processor/batch/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	DaoProvider        database.DaoProvider
	ProcessJobProvider queue.JobProvider[processor.MessageParams]
	BatchJobProvider   queue.JobProvider[batch.MessageParams]
}

var (
	Ref = ref.Root.MustSub("batcher")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides a queue worker for batch torrent processing jobs"),
		builder.WithActivation[deps](plugin.ActivationAlways),
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
						append(
							[]model.QueueJobOption{model.QueueJobMaxRetries(2)},
							options...)...,
					)
				}
			}),
		),
		builder.WithQueueHandler(
			func(deps deps) handler.Handler {
				return handler.New(
					Ref.String(),
					batch_queue.New(
						deps.DaoProvider,
						deps.ProcessJobProvider,
						deps.BatchJobProvider,
					),
					handler.JobTimeout(time.Second*60*10),
					handler.Concurrency(1),
				)
			},
		),
	)
)
