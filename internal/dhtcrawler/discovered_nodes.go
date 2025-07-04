package dhtcrawler

import (
	"context"
	"net/netip"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
)

// NewDiscoveredNodes creates the channel for discovered nodes.
// It receives nodes discovered by the crawler, as well as nodes from incoming requests to the DHT server.
// It is provided as a separate service to avoid a circular dependency with the DHT server.
func NewDiscoveredNodes(config Config) concurrency.BatchingChannel[ktable.Node] {
	return concurrency.NewBatchingChannel[ktable.Node](
		int(100*config.ScalingFactor),
		10,
		time.Second/100,
	)
}

func (cr *crawler) runDiscoveredNodes(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ps := <-cr.discoveredNodes.Out():
			addrs := make([]netip.Addr, 0, 1)

			m := make(map[string]ktable.Node, 1)
			for _, p := range ps {
				if _, ok := m[p.Addr().Addr().String()]; !ok {
					m[p.Addr().Addr().String()] = p
					addrs = append(addrs, p.Addr().Addr())
				}
			}
			// for any discovered node not already in the routing table,
			// we will block until it can be sent to any one of the pipeline channels.
			unknownAddrs := cr.kTable.FilterKnownAddrs(addrs)
			for _, addr := range unknownAddrs {
				p := m[addr.String()]
				select {
				case <-ctx.Done():
					return ctx.Err()
				case cr.nodesForFindNode.In() <- p:
				case cr.nodesForSampleInfoHashes.In() <- p:
				case cr.nodesForPing.In() <- p:
				}
			}
		}
	}
}
