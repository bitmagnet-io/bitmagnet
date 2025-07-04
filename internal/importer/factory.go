package importer

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
)

func New(daoProvider database.DaoTransactionProvider) Importer {
	return importer{
		daoProvider: daoProvider,
		bufferSize:  100,
		maxWaitTime: time.Second,
	}
}
