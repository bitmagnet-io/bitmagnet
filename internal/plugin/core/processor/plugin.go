package processor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/info_hash_blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	Indexer processor.Processor
}

var (
	Ref = ref.Root.MustSub("processor")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides the processor service that classifies and indexes torrents"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithDependencies[deps](
			classifier.Ref,
			info_hash_blocker.Ref,
			persister.Ref,
		),
		builder.WithFxOption[deps](
			fx.Provide(
				processor.New,
				func() queue.JobProvider[processor.MessageParams] {
					return func(msg processor.MessageParams, options ...model.QueueJobOption) (model.QueueJob, error) {
						return model.NewQueueJob(
							Ref.String(),
							msg,
							append(
								[]model.QueueJobOption{model.QueueJobMaxRetries(2)},
								options...)...,
						)
					}
				},
			),
		),
		builder.WithQueueHandler(
			func(deps deps) handler.Handler {
				return handler.New(
					Ref.String(),
					func(job model.QueueJob) runner.Runner {
						return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
							msg := &processor.MessageParams{}
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
