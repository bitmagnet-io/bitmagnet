package blockingfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/blocking"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		blocking.Namespace,
		fx.Provide(
			blocking.New,
			fx.Annotate(
				func(b blocking.Blocker) registry.Option {
					return registry.WithWorker(
						blocking.Namespace,
						b.Runner,
						worker.WithDependencies(database.Namespace),
					)
				},
				fx.ResultTags(`group:"worker_options"`),
			),
		),
	)
}
