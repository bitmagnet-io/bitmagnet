package dht_crawler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/channel"
	workers_metrics "github.com/bitmagnet-io/bitmagnet/internal/workers/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/periodic"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

func newPingWorker(
	cl client.Client,
	kTable ktable.Table,
	size atomic.Reader[PingConcurrency],
	oldPeerThreshold time.Duration,
	metrics *metrics.Component,
) channel.Worker[ktable.Node] {
	return channel.NewWorker(
		func(ctx context.Context, node ktable.Node) error {
			if node.Dropped() || node.Time().After(time.Now().Add(-oldPeerThreshold)) {
				// Either the node was already dropped or it succeeded after being added to the channel.
				// In either case we can continue.
				return nil
			}

			res, err := cl.Ping(ctx, node.Addr())

			var nodeID protocol.ID

			if err == nil {
				nodeID = res.ID
				if !node.ID().IsZero() && node.ID() != nodeID {
					nodeID = node.ID()
					err = errors.New("node responded with a mismatching ID")
				}
			}

			if err != nil {
				kTable.BatchCommand(ktable.DropNode{
					ID:     nodeID,
					Reason: fmt.Errorf("failed to respond to ping: %w", err),
				})
			} else {
				kTable.BatchCommand(ktable.PutNode{
					ID:      nodeID,
					Addr:    node.Addr(),
					Options: []ktable.NodeOption{ktable.NodeResponded()},
				},
				)
			}

			return nil
		},
		channel.WithAtomicSize[ktable.Node](
			atomic.MapIntish[PingConcurrency, int](size),
		),
		channel.WithQuickShutdown[ktable.Node](),
		channel.WithMetricsAdapter[ktable.Node](
			workers_metrics.MustNew(metrics.MustSub(dht.QPing)),
		),
	)
}

func newOldNodesWorker(
	kTable ktable.Table,
	interval time.Duration,
	oldPeerThreshold time.Duration,
	pingAdder channel.Adder[ktable.Node],
) runner.Provider {
	return periodic.New(
		interval,
		func(ctx context.Context) error {
			return pingAdder.Add(ctx, kTable.GetOldestNodes(time.Now().Add(-oldPeerThreshold), 0)...)
		},
		periodic.WithQuickShutdown(),
	)
}
