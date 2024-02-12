package postgres

import (
	"context"
	"database/sql"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"
	"sync"
)

type Params struct {
	fx.In
	Config Config
}

type Result struct {
	fx.Out
	PgxPool     lazy.Lazy[*pgxpool.Pool]
	SqlDB       lazy.Lazy[*sql.DB]
	PgxPoolWait *sync.WaitGroup `name:"pgx_pool_wait"`
	AppHook     fx.Hook         `group:"app_hooks"`
}

func New(p Params) (Result, error) {
	stopped := make(chan struct{})
	waitGroup := &sync.WaitGroup{}
	lazyPool := lazy.New(func() (*pgxpool.Pool, error) {
		ctx, cancel := context.WithCancel(context.Background())
		pl, err := pgxpool.New(ctx, p.Config.DSN())
		go func() {
			<-stopped
			// wait for services to be finished with the pool before closing
			waitGroup.Wait()
			cancel()
			pl.Close()
		}()
		return pl, err
	})
	return Result{
		PgxPool: lazyPool,
		SqlDB: lazy.New(func() (*sql.DB, error) {
			pool, err := lazyPool.Get()
			if err != nil {
				return nil, err
			}
			return stdlib.OpenDBFromPool(pool), nil
		}),
		PgxPoolWait: waitGroup,
		AppHook: fx.Hook{
			OnStop: func(context.Context) error {
				close(stopped)
				return nil
			},
		},
	}, nil
}
