package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Params struct {
	fx.In
	Config Config
	Logger *zap.SugaredLogger
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
		pl, plErr := pgxpool.New(ctx, p.Config.CreateDSN())
		if plErr != nil {
			cancel()
			return nil, plErr
		}
		if pingErr := waitForPing(ctx, p.Logger, pl); pingErr != nil {
			cancel()
			return nil, pingErr
		}
		go func() {
			<-stopped
			// wait for services to be finished with the pool before closing
			waitGroup.Wait()
			cancel()
			pl.Close()
		}()
		return pl, nil
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

func waitForPing(ctx context.Context, logger *zap.SugaredLogger, pool *pgxpool.Pool) error {
	i := 0
	var err error
	for {
		if ctx.Err() != nil {
			err = ctx.Err()
			break
		}
		err = pool.Ping(ctx)
		if err == nil {
			return nil
		}
		i++
		if i > 10 {
			break
		}
		select {
		case <-ctx.Done():
			break
		case <-time.After(time.Second):
			logger.Warnw("failed to ping database, retrying...", "error", err)
			break
		}
	}
	return fmt.Errorf("timed out waiting for ping: %w", err)
}
