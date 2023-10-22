package dhtcrawler

import (
	"context"
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
				_ = c.peersForFindNode.In(ctx, peer{addr: addr})
			}
		}
		interval = c.reseedBootstrapNodesInterval
	}
}
