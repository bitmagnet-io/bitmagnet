package dhtcrawler

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"time"
)

func (c *crawler) findNode(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.findNodesInterval):
			findNodesBatchSize := c.sampleInfoHashesShortfall.Get()
			if findNodesBatchSize < 1 {
				findNodesBatchSize = 4
			}
			peers := c.kTable.GetOldestPeers(time.Now().Add(-(5 * time.Second)), findNodesBatchSize)
			for _, p := range peers {
				_ = c.peersForFindNode.In(ctx, p)
			}
		}
	}
}

func (c *crawler) awaitPeersForFindNode(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case p := <-c.peersForFindNode.Out():
			if err := c.findNodeSemaphore.Acquire(ctx, 1); err != nil {
				break
			}
			go (func() {
				defer c.findNodeSemaphore.Release(1)
				id := c.soughtPeerId.Get()
				args := dht.MsgArgs{
					Target: id,
				}
				res, err := c.server.Query(ctx, p.Addr(), dht.QFindNode, args)
				if err != nil {
					c.kTable.BatchCommand(ktable.DropPeer{ID: p.ID(), Reason: fmt.Errorf("find_node failed: %w", err)})
				} else {
					if p.IsSampleInfoHashesCandidate() {
						_ = c.peersForSampleInfoHashes.TryIn(p)
					}
					c.kTable.BatchCommand(ktable.PutPeer{ID: p.ID(), Addr: p.Addr(), Options: []ktable.PeerOption{ktable.PeerResponded()}})
					for _, n := range res.Msg.R.Nodes {
						c.discoveredPeers.InContext(ctx, peer{
							id:   n.ID,
							addr: n.Addr.ToAddrPort(),
						})
					}
				}
			})()
		}
	}
}
