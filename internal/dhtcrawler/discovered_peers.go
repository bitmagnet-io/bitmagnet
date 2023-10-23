package dhtcrawler

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"net/netip"
)

func (c *crawler) awaitDiscoveredPeers(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ps := <-c.discoveredPeers.Out():
			m := make(map[string]ktable.Peer, len(ps))
			addrs := make([]netip.Addr, 0, len(ps))
			for _, p := range ps {
				m[p.Addr().Addr().String()] = p
				addrs = append(addrs, p.Addr().Addr())
			}
			unknownAddrs := c.kTable.FilterKnownAddrs(addrs)
			for _, addr := range unknownAddrs {
				if p, ok := m[addr.String()]; ok {
					if err := c.peersForSampleInfoHashes.TryIn(p); err != nil {
						_ = c.peersForFindNode.In(ctx, p)
					}
				}
			}
		}
	}
}
