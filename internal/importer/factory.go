package importer

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"go.uber.org/fx"
	"time"
)

type Params struct {
	fx.In
	Dao lazy.Lazy[*dao.Query]
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
			return importer{
				dao:         d,
				bufferSize:  100,
				maxWaitTime: 500 * time.Millisecond,
			}, nil
		}),
	}
}
