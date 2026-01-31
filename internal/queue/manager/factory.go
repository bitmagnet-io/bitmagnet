package manager

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/batch"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
)

func New(
	db database.GormDBProvider,
	batchJobProvider queue.JobProvider[batch.MessageParams],
) Manager {
	return manager{
		db:               db,
		batchJobProvider: batchJobProvider,
	}
}
