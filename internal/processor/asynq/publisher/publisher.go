package publisher

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/producer"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/publisher"
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Client   lazy.Lazy[*asynq.Client]
	Producer producer.Producer[processor.MessageParams]
}

type Result struct {
	fx.Out
	Publisher lazy.Lazy[publisher.Publisher[processor.MessageParams]]
}

func New(p Params) Result {
	return Result{
		Publisher: lazy.New(func() (publisher.Publisher[processor.MessageParams], error) {
			client, err := p.Client.Get()
			if err != nil {
				return nil, err
			}
			return publisher.New[processor.MessageParams](client, p.Producer), nil
		}),
	}
}
