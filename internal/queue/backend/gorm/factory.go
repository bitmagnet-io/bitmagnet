package gorm

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/exp/slog"
)

type Params struct {
	fx.In
	Query   lazy.Lazy[*dao.Query]
	PgxPool lazy.Lazy[*pgxpool.Pool]
	Logger  *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Queue   lazy.Lazy[queue.Queue]
	AppHook fx.Hook `group:"app_hooks"`
}

func New(p Params) Result {
	q := lazy.New(func() (queue.Queue, error) {
		pool, err := p.PgxPool.Get()
		if err != nil {
			return nil, err
		}
		query, err := p.Query.Get()
		if err != nil {
			return nil, err
		}
		cfg := queue.NewConfig()
		cfg.IdleTransactionTimeout = queue.DefaultIdleTxTimeout
		cfg.PGConnectionTimeout = DefaultConnectionTimeout

		ctx, cancel := context.WithCancel(context.Background())

		b := &backend{
			query:          query,
			pool:           pool,
			cancelFuncs:    []context.CancelFunc{cancel},
			config:         cfg,
			handlers:       make(map[string]handler.Handler),
			newQueues:      make(chan string),
			readyQueues:    make(chan string),
			listenCancelCh: make(chan context.CancelFunc, 1),
			logger:         p.Logger.Named("queue"),
		}

		//backend.listenerConn, err = backend.newListenerConn(ctx)
		if err != nil {
			b.logger.Error("unable to initialize listener connection", slog.Any("error", err))
		}

		// monitor handlers for changes and LISTEN when new queues are added
		go b.newQueueMonitor(ctx)

		return b, nil
	})
	return Result{
		Queue: q,
		AppHook: fx.Hook{
			OnStop: func(ctx context.Context) error {
				return q.IfInitialized(func(q queue.Queue) error {
					q.Shutdown(ctx)
					return nil
				})
			},
		},
	}
}
