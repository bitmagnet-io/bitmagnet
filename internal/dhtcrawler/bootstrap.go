package dhtcrawler

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"net"
	"time"
)

func (c *crawler) reseedBootstrapNodes(ctx context.Context) {
	interval := time.Duration(0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			for _, strAddr := range c.bootstrapNodes {
				addr, err := net.ResolveUDPAddr("udp", strAddr)
				if err != nil {
					c.logger.Warnf("failed to resolve bootstrap node address: %s", err)
					continue
				}
				select {
				case <-ctx.Done():
					return
				case c.nodesForPing.In() <- ktable.NewNode(ktable.ID{}, addr.AddrPort()):
					continue
				}
			}
		}
		interval = c.reseedBootstrapNodesInterval
	}
}
