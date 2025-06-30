package dhtcrawler

import (
	"context"
	"net/netip"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

type DiscoveredNodesParams struct {
	fx.In
	Config Config
}

type DiscoveredNodesResult struct {
	fx.Out
	DiscoveredNodes concurrency.BatchingChannel[ktable.Node] `name:"dht_discovered_nodes"`
}

// NewDiscoveredNodes creates the channel for discovered nodes.
// It receives nodes discovered by the crawler, as well as nodes from incoming requests to the DHT server.
// It is provided as a separate service to avoid a circular dependency with the DHT server.
func NewDiscoveredNodes(params DiscoveredNodesParams) DiscoveredNodesResult {
	return DiscoveredNodesResult{
		DiscoveredNodes: concurrency.NewBatchingChannel[ktable.Node](
			int(100*params.Config.ScalingFactor), 10, time.Second/100),
	}
}

func (c *crawler) runDiscoveredNodes(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ps := <-c.discoveredNodes.Out():
			addrs := make([]netip.Addr, 0, 1)

			m := make(map[string]ktable.Node, 1)
			for _, p := range ps {
				if _, ok := m[p.Addr().Addr().String()]; !ok {
					m[p.Addr().Addr().String()] = p
					addrs = append(addrs, p.Addr().Addr())
				}
			}

			// For any newly discovered node, we will block until it can be
			// sent to any one of the pipeline channels.
			for _, addr := range addrs {
				p := m[addr.String()]

				var result string

				if c.ignoreNodes.test(p.ID()) {
					result = "ignored"
				} else {
					select {
					case <-ctx.Done():
						return
					case c.nodesForFindNode.In() <- p:
						result = "find_node"
					case c.nodesForSampleInfoHashes.In() <- p:
						c.ignoreNodes.add(p.ID())
						result = "sample_infohashes"
					case c.nodesForPing.In() <- p:
						result = "ping"
					}
				}

				c.discoveredNodesTotal.With(prometheus.Labels{"result": result}).Inc()
			}
		}
	}
}
