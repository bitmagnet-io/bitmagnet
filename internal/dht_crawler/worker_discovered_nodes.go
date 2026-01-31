package dht_crawler

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/deduplicator"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/channel"
)

func newDiscoveredNodesWorker(
	adders func() []channel.Adder[ktable.Node],
	timeout time.Duration,
	size int,
) channel.Worker[ktable.Node] {
	seenNodes := deduplicator.New[string](100_000, time.Hour)

	return channel.NewWorker(
		func(ctx context.Context, node ktable.Node) error {
			// Skip recently seen nodes; bootstrap nodes have a zero value ID and are not skipped.
			if !node.ID().IsZero() && !seenNodes.Add(node.Addr().String()) {
				return nil
			}

			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			for _, adder := range adders() {
				go func(adder channel.Adder[ktable.Node]) {
					if adder.Add(ctx, node) == nil {
						cancel()
					}
				}(adder)
			}

			<-ctx.Done()

			return nil
		},
		channel.WithSize[ktable.Node](size),
		channel.WithQuickShutdown[ktable.Node](),
	)
}
