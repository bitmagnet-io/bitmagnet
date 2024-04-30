package processorfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/queue"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"processor",
		fx.Provide(
			processor.New,
			queue.New,
		),
	)
}
