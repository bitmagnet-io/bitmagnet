package dhtcrawler

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"time"
)

func (c *crawler) runPing(ctx context.Context) {
	args := dht.MsgArgs{
		ID: c.kTable.Origin(),
	}
	_ = c.nodesForPing.Run(ctx, func(n ktable.Node) {
		if n.Dropped() || n.Time().After(time.Now().Add(-c.oldPeerThreshold)) {
			// Either the node was already dropped or it succeeded after being added to the channel.
			// In either case we can continue.
			return
		}
		res, err := c.server.Query(ctx, n.Addr(), dht.QPing, args)
		var nodeID protocol.ID
		if err == nil {
			nodeID = res.Msg.R.ID
			if !n.ID().IsZero() && n.ID() != nodeID {
				nodeID = n.ID()
				err = errors.New("node responded with a mismatching ID")
			}
		}
		if err != nil {
			c.kTable.BatchCommand(ktable.DropNode{ID: nodeID, Reason: fmt.Errorf("failed to respond to ping: %w", err)})
		} else {
			c.kTable.BatchCommand(ktable.PutNode{ID: nodeID, Addr: n.Addr(), Options: []ktable.NodeOption{ktable.NodeResponded()}})
		}
	})
}

func (c *crawler) getOldNodes(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.getOldestNodesInterval):
			for _, p := range c.kTable.GetOldestNodes(time.Now().Add(-c.oldPeerThreshold), 0) {
				select {
				case <-ctx.Done():
					return
				case c.nodesForPing.In() <- p:
					continue
				}
			}
		}
	}
}
