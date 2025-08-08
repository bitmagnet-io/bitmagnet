package dhtcrawler

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	size int,
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
		channel.WithSize[ktable.Node](size),
		channel.WithQuickShutdown[ktable.Node](),
		channel.WithMetricsAdapter[ktable.Node](
			workers_metrics.MustNew(metrics.MustSub(dht.QPing)),
		),
	)
}

// func (cr *crawler) runPing(ctx context.Context) error {
// 	return cr.nodesForPing.Run(ctx, func(n ktable.Node) {
// 		if n.Dropped() || n.Time().After(time.Now().Add(-cr.oldPeerThreshold)) {
// 			// Either the node was already dropped or it succeeded after being added to the channel.
// 			// In either case we can continue.
// 			return
// 		}

// 		res, err := cr.client.Ping(ctx, n.Addr())

// 		var nodeID protocol.ID

// 		if err == nil {
// 			nodeID = res.ID
// 			if !n.ID().IsZero() && n.ID() != nodeID {
// 				nodeID = n.ID()
// 				err = errors.New("node responded with a mismatching ID")
// 			}
// 		}

// 		if err != nil {
// 			cr.kTable.BatchCommand(ktable.DropNode{
// 				ID:     nodeID,
// 				Reason: fmt.Errorf("failed to respond to ping: %w", err),
// 			})
// 		} else {
// 			cr.kTable.BatchCommand(ktable.PutNode{
// 				ID:      nodeID,
// 				Addr:    n.Addr(),
// 				Options: []ktable.NodeOption{ktable.NodeResponded()},
// 			},
// 			)
// 		}
// 	})
// }

// // getOldNodes periodically adds the oldest nodes from the routing table to the nodesForPing channel,
// // so they can be pruned from the routing table if no longer responsive.
// func (cr *crawler) getOldNodes(ctx context.Context) error {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		case <-time.After(cr.getOldestNodesInterval):
// 			for _, p := range cr.kTable.GetOldestNodes(time.Now().Add(-cr.oldPeerThreshold), 0) {
// 				select {
// 				case <-ctx.Done():
// 					return ctx.Err()
// 				case cr.nodesForPing.In() <- p:
// 					continue
// 				}
// 			}
// 		}
// 	}
// }

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
