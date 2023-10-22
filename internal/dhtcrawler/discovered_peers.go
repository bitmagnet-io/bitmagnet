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
			//// before pinging a discovered peer we'll check if we already know about it;
			//// we may or may not know the peer ID, and that will determine how we'll carry out this check
			//if p.ID().IsZero() {
			//	// we don't know the peer ID, so we'll do a reverse lookup on the address
			//	if c.kTable.HasAddr(p.Addr().Addr()) {
			//		break
			//	}
			//} else {
			//	// we know the peer ID, so we can check if it's already in the routing table
			//	if c.kTable.HasPeer(p.ID()) {
			//		break
			//	}
			//}
			//if err := c.peersForSampleInfoHashes.TryIn(p); err != nil {
			//	_ = c.peersForFindNode.In(ctx, p)
			//}
		}
	}
}
