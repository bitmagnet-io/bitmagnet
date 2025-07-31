package dhtcrawlerhealthcheck

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
)

func New(
	name string,
	isActive func() bool,
	dhtServerLastResponses *concurrency.AtomicValue[server.LastResponses],
	// isActive *concurrency.AtomicValue[bool],
) health.CheckerOption {
	return health.WithPeriodicCheck(
		time.Second*10,
		time.Second*1,
		NewCheck(name, isActive, dhtServerLastResponses),
	)
}
