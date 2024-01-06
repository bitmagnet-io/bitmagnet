package consumer

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/message"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/consumer"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Classifier lazy.Lazy[classifier.Classifier]
}

type Result struct {
	fx.Out
	Consumer lazy.Lazy[consumer.Consumer] `group:"queue_consumers"`
}

func New(p Params) Result {
	return Result{
		Consumer: lazy.New(func() (consumer.Consumer, error) {
			cl, err := p.Classifier.Get()
			if err != nil {
				return nil, err
			}
			return consumer.New[message.ClassifyTorrentPayload](
				message.ClassifyTorrentTypename,
				cns{
					cl,
				},
			), nil
		}),
	}
}

type cns struct {
	c classifier.Classifier
}

func (c cns) Handle(ctx context.Context, msg message.ClassifyTorrentPayload) error {
	return c.c.Classify(ctx, msg.InfoHashes...)
}
