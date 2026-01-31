package persister

import (
	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/indexer"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"go.uber.org/zap"
)

func New(
	maxSize *atomic.Value[MaxSize],
	maxWait *atomic.Value[MaxWait],
	daoProvider database.DaoTransactionProvider,
	blockerBlocker blocker.Blocker,
	logger *zap.Logger,
	metrics *metrics.Component,
	indexers []indexer.Indexer,
) Persister {
	flusher := iflusher(newFlusherMetrics(
		&flusher{
			daoProvider: daoProvider,
			blocker:     blockerBlocker,
			sem:         make(chan struct{}, 1),
		},
		metrics,
	))

	if len(indexers) > 0 {
		flusher = &flusherIndexer{
			iflusher: flusher,
			indexer:  indexer.Multi(indexers),
		}
	}

	return &worker{
		iflusher: flusher,
		shutdown: make(chan struct{}),
		maxSize:  maxSize,
		maxWait:  maxWait,
		logger:   logger,
	}
}
