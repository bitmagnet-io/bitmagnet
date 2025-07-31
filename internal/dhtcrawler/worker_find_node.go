package dhtcrawler

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/channel"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/periodic"
)

// func (cr *crawler) getNodesForFindNode(ctx context.Context) error {
// 	for {
// 		peers := cr.kTable.GetOldestNodes(time.Now().Add(-(5 * time.Second)), 10)
// 		for _, p := range peers {
// 			select {
// 			case <-ctx.Done():
// 				return ctx.Err()
// 			case cr.nodesForFindNode.In() <- p:
// 				continue
// 			}
// 		}

// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		case <-time.After(time.Second):
// 		}
// 	}
// }

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
	soughtNodeID *concurrency.AtomicValue[protocol.ID],
	size int,
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
		channel.WithSize[ktable.Node](size),
	)
}

// func (cr *crawler) runFindNode(ctx context.Context) error {
// 	return cr.nodesForFindNode.Run(ctx, func(p ktable.Node) {
// 		res, err := cr.client.FindNode(ctx, p.Addr(), cr.soughtNodeID.Get())
// 		if err != nil {
// 			cr.kTable.BatchCommand(ktable.DropNode{
// 				ID:     p.ID(),
// 				Reason: fmt.Errorf("find_node failed: %w", err),
// 			})
// 		} else {
// 			cr.kTable.BatchCommand(ktable.PutNode{
// 				ID:      p.ID(),
// 				Addr:    p.Addr(),
// 				Options: []ktable.NodeOption{ktable.NodeResponded()},
// 			})
// 			// block this channel until all nodes can be added to the discoveredNodes channel
// 			for _, n := range res.Nodes {
// 				select {
// 				case <-ctx.Done():
// 					return
// 				case cr.discoveredNodes.In() <- ktable.NewNode(n.ID, n.Addr):
// 					continue
// 				}
// 			}
// 		}
// 	})
// }
