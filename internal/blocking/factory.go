package blocking

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"go.uber.org/fx"
	"time"
)

type Params struct {
	fx.In
	Dao *dao.Query
}

type Result struct {
	fx.Out
	Manager Manager
	AppHook fx.Hook `group:"app_hooks"`
}

func New(params Params) Result {
	m := &manager{
		dao:           params.Dao,
		buffer:        make(map[protocol.ID]struct{}, 1000),
		maxBufferSize: 1000,
		maxFlushWait:  time.Minute * 5,
	}
	return Result{
		Manager: m,
		AppHook: fx.Hook{
			OnStop: func(ctx context.Context) error {
				return m.Flush(ctx)
			},
		},
	}
}
