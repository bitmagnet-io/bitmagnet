package processorfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/queue/decorator"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/queue/publisher"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/worker"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"processor",
		fx.Provide(
			processor.New,
			decorator.New,
			publisher.New,
			worker.New,
		),
	)
}
