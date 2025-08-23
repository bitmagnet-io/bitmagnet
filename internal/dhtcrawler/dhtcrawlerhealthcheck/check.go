package dhtcrawlerhealthcheck

import (
	"context"
	"errors"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
)

func NewCheck(
	name string,
	isActive func() bool,
	//isActive *concurrency.Value[bool],
	lastResponses *atomic.Value[server.LastResponses],
) health.Check {
	return health.Check{
		Name:     name,
		IsActive: isActive,
		Timeout:  time.Second,
		Check: func(context.Context) error {
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
