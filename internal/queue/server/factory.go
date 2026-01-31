package server

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"go.uber.org/zap"
)

func New(
	daoProvider database.DaoTransactionProvider,
	logger *zap.Logger,
	handlers ...handler.Handler,
) runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
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
