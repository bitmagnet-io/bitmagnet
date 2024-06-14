package processorfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	batchqueue "github.com/bitmagnet-io/bitmagnet/internal/processor/batch/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/hook_0_9_0"
	processorqueue "github.com/bitmagnet-io/bitmagnet/internal/processor/queue"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"processor",
		configfx.NewConfigModule[processor.Config]("processor", processor.NewDefaultConfig()),
		fx.Provide(
			processor.New,
			processorqueue.New,
			batchqueue.New,
		),
		fx.Decorate(
			hook_0_9_0.NewDecorator,
		),
	)
}
