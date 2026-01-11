package dht_crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/channel"
	workers_metrics "github.com/bitmagnet-io/bitmagnet/internal/workers/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/periodic"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

func newGetNodesForFindNodeWorker(
	kTable ktable.Table,
	findNodesAdder channel.Adder[ktable.Node],
) runner.Provider {
	return periodic.New(
		time.Minute,
		func(ctx context.Context) error {
			nodes := kTable.GetOldestNodes(time.Now().Add(-(5 * time.Second)), 10)
			err := findNodesAdder.Add(ctx, nodes...)
			if err != nil {
				return err
			}

			return nil
		},
		periodic.WithQuickShutdown(),
	)
}

func newFindNodesWorker(
	cl client.Client,
	kTable ktable.Table,
	discoveredNodesAdder channel.Adder[ktable.Node],
	soughtNodeID *atomic.Value[protocol.ID],
	size atomic.Reader[FindNodesConcurrency],
	metrics *metrics.Component,
	bootstrap func(context.Context) error,
) channel.Worker[ktable.Node] {
	return channel.NewWorker(
		func(ctx context.Context, node ktable.Node) error {
			res, err := cl.FindNode(ctx, node.Addr(), soughtNodeID.Get())

			if err != nil {
				kTable.BatchCommand(ktable.DropNode{
					ID:     node.ID(),
					Reason: fmt.Errorf("find_node failed: %w", err),
				})

				return nil
			} else {
				kTable.BatchCommand(ktable.PutNode{
					ID:      node.ID(),
					Addr:    node.Addr(),
					Options: []ktable.NodeOption{ktable.NodeResponded()},
				})
			}

			return discoveredNodesAdder.Add(ctx, slice.Map(res.Nodes, func(info client.NodeInfo) ktable.Node {
				return ktable.NewNode(info.ID, info.Addr)
			})...)
		},
		channel.WithQuickShutdown[ktable.Node](),
		channel.WithAtomicSize[ktable.Node](
			atomic.MapIntish[FindNodesConcurrency, int](size),
		),
		channel.WithMetricsAdapter[ktable.Node](
			workers_metrics.MustNew(metrics.MustSub(dht.QFindNode)),
		),
		// If this worker becomes idle we should bootstrap again:
		channel.WithOnIdle[ktable.Node](bootstrap),
	)
}
