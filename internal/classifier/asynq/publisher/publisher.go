package publisher

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/message"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/producer"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/publisher"
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Client   lazy.Lazy[*asynq.Client]
	Producer producer.Producer[message.ClassifyTorrentPayload]
}

type Result struct {
	fx.Out
	Publisher lazy.Lazy[publisher.Publisher[message.ClassifyTorrentPayload]]
}

func New(p Params) Result {
	return Result{
		Publisher: lazy.New(func() (publisher.Publisher[message.ClassifyTorrentPayload], error) {
			client, err := p.Client.Get()
			if err != nil {
				return nil, err
			}
			return publisher.New[message.ClassifyTorrentPayload](client, p.Producer), nil
		}),
	}
}
