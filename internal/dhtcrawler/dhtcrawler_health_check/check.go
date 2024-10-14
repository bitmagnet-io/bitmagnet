package dhtcrawler_health_check

import (
	"context"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
	"time"
)

func NewCheck(
	dhtCrawlerActive *concurrency.AtomicValue[bool],
	lastResponses *concurrency.AtomicValue[server.LastResponses],
) health.Check {
	return health.Check{
		Name: "dht",
		IsActive: func() bool {
			return dhtCrawlerActive.Get()
		},
		Timeout: time.Second,
		Check: func(ctx context.Context) error {
			lr := lastResponses.Get()
			if lr.StartTime.IsZero() {
				return nil
			}
			now := time.Now()
			if lr.LastSuccess.IsZero() {
				if now.Sub(lr.StartTime) < 30*time.Second {
					return nil
				}
				return errors.New("no response within 30 seconds")
			}
			if now.Sub(lr.LastSuccess) > time.Minute {
				return errors.New("no successful responses within last minute")
			}
			return nil
		},
	}
}
