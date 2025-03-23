package blocking

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"sync"
	"time"
)

type Params struct {
	fx.In
	Pool        lazy.Lazy[*pgxpool.Pool]
	PgxPoolWait *sync.WaitGroup `name:"pgx_pool_wait"`
}

type Result struct {
	fx.Out
	Manager lazy.Lazy[Manager]
	AppHook fx.Hook `group:"app_hooks"`
}

func New(params Params) Result {
	lazyManager := lazy.New[Manager](func() (Manager, error) {
		pool, err := params.Pool.Get()
		if err != nil {
			return nil, err
		}
		params.PgxPoolWait.Add(1)
		return &manager{
			pool:          pool,
			buffer:        make(map[protocol.ID]struct{}, 1000),
			maxBufferSize: 1000,
			maxFlushWait:  time.Minute * 5,
		}, nil
	})
	return Result{
		Manager: lazyManager,
		AppHook: fx.Hook{
			OnStop: func(ctx context.Context) error {
				return lazyManager.IfInitialized(func(m Manager) error {
					defer params.PgxPoolWait.Done()
					return m.Flush(ctx)
				})
			},
		},
	}
}
