package dhtcrawler

import (
	"context"
	"net"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
)

func (cr *crawler) reseedBootstrapNodes(ctx context.Context) error {
	interval := time.Duration(0)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
			for _, strAddr := range cr.bootstrapNodes {
				addr, err := net.ResolveUDPAddr("udp", strAddr)
				if err != nil {
					cr.logger.Warnf("failed to resolve bootstrap node address: %s", err)
					continue
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case cr.nodesForPing.In() <- ktable.NewNode(ktable.ID{}, addr.AddrPort()):
					continue
				}
			}
		}

		interval = cr.reseedBootstrapNodesInterval
	}
}
