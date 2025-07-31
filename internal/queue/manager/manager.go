package manager

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/indexer/batch"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
)

type manager struct {
	batchJobProvider queue.JobProvider[batch.MessageParams]
	db               database.GormDBProvider
}
