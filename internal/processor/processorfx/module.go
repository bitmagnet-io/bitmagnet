package processorfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/asynq/consumer"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/asynq/decorator"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/asynq/producer"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/asynq/publisher"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"processor",
		fx.Provide(
			processor.New,
			consumer.New,
			decorator.New,
			producer.New,
			publisher.New,
		),
	)
}
