package server

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"go.uber.org/zap"
)

func New(
	daoProvider database.DaoTransactionProvider,
	logger *zap.SugaredLogger,
	lazyHandlers ...lazy.Lazy[handler.Handler],
) runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		handlers := make([]handler.Handler, 0, len(lazyHandlers))

		for _, lh := range lazyHandlers {
			h, err := lh.Get()
			if err != nil {
				return runner.NopShutdowner, err
			}

			handlers = append(handlers, h)
		}

		srv := server{
			daoProvider: daoProvider,
			handlers:    handlers,
			gcInterval:  time.Minute * 10,
			gcSemaphore: make(chan struct{}, 1),
			draining:    make(chan struct{}),
			logger:      logger,
		}

		return srv.Start(ctx, cancel)
	}
}
