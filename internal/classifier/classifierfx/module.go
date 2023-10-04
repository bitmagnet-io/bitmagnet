package classifierfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/consumer"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/producer"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/publisher"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/resolver/video"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/video/videofx"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"classifier",
		fx.Provide(
			classifier.New,
			consumer.New,
			producer.New,
			publisher.New,
			video.New,
		),
		videofx.New(),
	)
}
