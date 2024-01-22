package consumer

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/consumer"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Processor lazy.Lazy[processor.Processor]
}

type Result struct {
	fx.Out
	Consumer lazy.Lazy[consumer.Consumer] `group:"queue_consumers"`
}

func New(p Params) Result {
	return Result{
		Consumer: lazy.New(func() (consumer.Consumer, error) {
			pr, err := p.Processor.Get()
			if err != nil {
				return nil, err
			}
			return consumer.New[processor.MessageParams](
				processor.MessageName,
				cns{
					pr,
				},
			), nil
		}),
	}
}

type cns struct {
	p processor.Processor
}

func (c cns) Handle(ctx context.Context, params processor.MessageParams) error {
	return c.p.Process(ctx, params)
}
