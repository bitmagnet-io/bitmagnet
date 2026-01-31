package queuemetrics

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
)

func New(dbProvider database.GormDBProvider) Client {
	return client{dbProvider: dbProvider}
}
