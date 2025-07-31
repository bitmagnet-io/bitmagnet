package blocker

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

func New(pool database.PoolProvider) Blocker {
	return &manager{
		pool:          pool,
		buffer:        make(map[protocol.ID]struct{}, 1000),
		maxBufferSize: 1000,
		maxFlushWait:  time.Minute * 5,
	}
}
