package dhtcrawler

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"time"
)

func (c *crawler) getNodesForSampleInfoHashes(ctx context.Context) {
	for {
		peers := c.kTable.GetNodesForSampleInfoHashes(60)
		for _, p := range peers {
			select {
			case <-ctx.Done():
				return
			case c.nodesForSampleInfoHashes.In() <- p:
				continue
			}
		}
	}
}

func (c *crawler) runSampleInfoHashes(ctx context.Context) {
	_ = c.nodesForSampleInfoHashes.Run(ctx, func(n ktable.Node) {
		if !n.IsSampleInfoHashesCandidate() {
			return
		}
		target := c.soughtNodeID.Get()
		args := dht.MsgArgs{
			ID:     c.kTable.Origin(),
			Target: target,
		}
		res, err := c.server.Query(ctx, n.Addr(), dht.QSampleInfohashes, args)
		if err != nil {
			c.kTable.BatchCommand(ktable.DropNode{ID: n.ID(), Reason: fmt.Errorf("sample_infohashes failed: %w", err)})
			return
		}
		var discoveredHashes []nodeHasPeersForHash
		if res.Msg.R.Samples != nil {
			for _, s := range *res.Msg.R.Samples {
				toa := c.ignoreHashes.testOrAdd(s)
				if !toa {
					discoveredHashes = append(discoveredHashes, nodeHasPeersForHash{
						infoHash: s,
						node:     n.Addr(),
					})
				}
			}
		}
		for _, h := range discoveredHashes {
			select {
			case <-ctx.Done():
				return
			case c.infoHashTriage.In() <- h:
				continue
			}
		}
		totalNum := 0
		interval := 0
		if res.Msg.R.Num != nil {
			totalNum = int(*res.Msg.R.Num)
		}
		if res.Msg.R.Interval != nil {
			interval = int(*res.Msg.R.Interval)
		}
		// most peers request a 6 hour backoff time(!)
		// if we're still discovering info hashes from them then let's set a respectful interval instead
		if len(discoveredHashes) > 0 && interval > 300 {
			interval = 60
		}
		c.kTable.BatchCommand(ktable.PutNode{ID: n.ID(), Addr: n.Addr(), Options: []ktable.NodeOption{
			ktable.NodeResponded(),
			ktable.NodeBep51Support(true),
			ktable.NodeSampleInfoHashesRes(
				len(discoveredHashes),
				totalNum,
				time.Now().Add(time.Duration(interval)*time.Second),
			),
		}})
		go func() {
			timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			for _, n := range res.Msg.R.Nodes {
				select {
				case <-timeoutCtx.Done():
					return
				case c.discoveredNodes.In() <- ktable.NewNode(n.ID, n.Addr.ToAddrPort()):
					continue
				}
			}
		}()
	})
}
