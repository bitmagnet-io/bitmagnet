package importer

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
)

func New(
	daoProvider database.DaoTransactionProvider,
	queueJobProvider queue.JobProvider[processor.MessageParams],
) Importer {
	return importer{
		daoProvider:      daoProvider,
		bufferSize:       100,
		maxWaitTime:      time.Second,
		queueJobProvider: queueJobProvider,
	}
}
