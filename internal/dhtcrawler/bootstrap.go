package dhtcrawler

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"time"
)

func (c *crawler) reseedBootstrapNodes(ctx context.Context) {
	interval := time.Duration(0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			for _, addr := range c.bootstrapNodes {
				select {
				case <-ctx.Done():
					return
				case c.nodesForPing.In() <- ktable.NewNode(ktable.ID{}, addr):
					continue
				}
			}
		}
		interval = c.reseedBootstrapNodesInterval
	}
}
