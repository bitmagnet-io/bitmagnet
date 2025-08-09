package persister

import (
	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"go.uber.org/zap"
)

func New(
	config Config,
	daoProvider database.DaoTransactionProvider,
	blockerBlocker blocker.Blocker,
	logger *zap.Logger,
	metrics *metrics.Component,
) Persister {
	return &worker{
		iflusher: newFlusherMetrics(
			&flusher{
				daoProvider: daoProvider,
				blocker:     blockerBlocker,
				sem:         make(chan struct{}, 1),
			},
			metrics,
		),
		shutdown: make(chan struct{}),
		maxSize:  config.MaxSize,
		maxWait:  config.MaxWait,
		logger:   logger,
	}
}
