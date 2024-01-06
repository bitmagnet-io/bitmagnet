package importer

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/message"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/publisher"
	"go.uber.org/fx"
	"time"
)

type Params struct {
	fx.In
	Dao               lazy.Lazy[*dao.Query]
	ClassifyPublisher lazy.Lazy[publisher.Publisher[message.ClassifyTorrentPayload]]
}

type Result struct {
	fx.Out
	Importer lazy.Lazy[Importer]
}

func New(p Params) Result {
	return Result{
		Importer: lazy.New(func() (Importer, error) {
			d, err := p.Dao.Get()
			if err != nil {
				return nil, err
			}
			cp, err := p.ClassifyPublisher.Get()
			if err != nil {
				return nil, err
			}
			return importer{
				dao:               d,
				classifyPublisher: cp,
				bufferSize:        100,
				maxWaitTime:       500 * time.Millisecond,
			}, nil
		}),
	}
}
