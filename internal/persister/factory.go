package persister

import (
	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"go.uber.org/zap"
)

func New(
	config Config,
	daoProvider database.DaoTransactionProvider,
	blockerBlocker blocker.Blocker,
	logger *zap.Logger,
) Persister {
	return &worker{
		flusher: &flusher{
			daoProvider: daoProvider,
			blocker:     blockerBlocker,
			sem:         make(chan struct{}, 1),
		},
		shutdown: make(chan struct{}),
		maxSize:  config.MaxSize,
		maxWait:  config.MaxWait,
		logger:   logger,
	}
}
