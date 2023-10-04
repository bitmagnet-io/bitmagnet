package consumer

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/message"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/consumer"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Classifier classifier.Classifier
}

type Result struct {
	fx.Out
	Consumer consumer.Consumer `group:"queue_consumers"`
}

func New(p Params) (Result, error) {
	return Result{
		Consumer: consumer.New[message.ClassifyTorrentPayload](
			message.ClassifyTorrentTypename,
			cns{
				p.Classifier,
			},
		),
	}, nil
}

type cns struct {
	c classifier.Classifier
}

func (c cns) Handle(ctx context.Context, msg message.ClassifyTorrentPayload) error {
	return c.c.Classify(ctx, msg.InfoHashes...)
}
