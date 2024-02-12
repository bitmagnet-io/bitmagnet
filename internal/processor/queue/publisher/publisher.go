package publisher

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/publisher"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Client lazy.Lazy[queue.Queue]
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
			return publisher.New[processor.MessageParams](client, processor.MessageName), nil
		}),
	}
}
