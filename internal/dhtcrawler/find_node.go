package dhtcrawler

import (
	"context"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
)

func (cr *crawler) getNodesForFindNode(ctx context.Context) error {
	for {
		peers := cr.kTable.GetOldestNodes(time.Now().Add(-(5 * time.Second)), 10)
		for _, p := range peers {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case cr.nodesForFindNode.In() <- p:
				continue
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
}

func (cr *crawler) runFindNode(ctx context.Context) error {
	return cr.nodesForFindNode.Run(ctx, func(p ktable.Node) {
		res, err := cr.client.FindNode(ctx, p.Addr(), cr.soughtNodeID.Get())
		if err != nil {
			cr.kTable.BatchCommand(ktable.DropNode{
				ID:     p.ID(),
				Reason: fmt.Errorf("find_node failed: %w", err),
			})
		} else {
			cr.kTable.BatchCommand(ktable.PutNode{
				ID:      p.ID(),
				Addr:    p.Addr(),
				Options: []ktable.NodeOption{ktable.NodeResponded()},
			})
			// block this channel until all nodes can be added to the discoveredNodes channel
			for _, n := range res.Nodes {
				select {
				case <-ctx.Done():
					return
				case cr.discoveredNodes.In() <- ktable.NewNode(n.ID, n.Addr):
					continue
				}
			}
		}
	})
}
