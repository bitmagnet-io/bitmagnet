package dhtcrawler

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"time"
)

func (c *crawler) sampleInfoHashes(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.sampleInfoHashesInterval):
			batchSize := (c.targetStagingSize - c.staging.count()) / 8
			if batchSize < 1 {
				c.sampleInfoHashesShortfall.Set(0)
				break
			}
			peers := c.kTable.GetPeersForSampleInfoHashes(batchSize)
			c.sampleInfoHashesShortfall.Set(batchSize - len(peers))
			for _, p := range peers {
				_ = c.peersForSampleInfoHashes.In(ctx, p)
			}
		}
	}
}

func (c *crawler) awaitPeersForSampleInfoHashes(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case p := <-c.peersForSampleInfoHashes.Out():
			if err := c.sampleInfoHashesSemaphore.Acquire(ctx, 1); err != nil {
				break
			}
			go (func() {
				defer c.sampleInfoHashesSemaphore.Release(1)
				if !p.IsSampleInfoHashesCandidate() {
					return
				}
				target := c.soughtPeerId.Get()
				args := dht.MsgArgs{
					ID:     c.kTable.Origin(),
					Target: target,
				}
				res, err := c.server.Query(ctx, p.Addr(), dht.QSampleInfohashes, args)
				var hashesToStage []infoHashWithPeer
				if err == nil {
					for _, n := range res.Msg.R.Nodes {
						_ = c.discoveredPeers.TryIn(peer{
							id:   n.ID,
							addr: n.Addr.ToAddrPort(),
						})
					}
					if res.Msg.R.Samples != nil {
						for _, s := range *res.Msg.R.Samples {
							hashesToStage = append(hashesToStage, infoHashWithPeer{
								infoHash: s,
								peer:     p.Addr(),
							})
						}
					}
					if len(hashesToStage) == 0 {
						err = errors.New("empty samples")
					}
				}
				if err != nil {
					c.kTable.BatchCommand(ktable.DropPeer{ID: p.ID(), Reason: fmt.Errorf("sample_infohashes failed: %w", err)})
					return
				}
				totalNum := 0
				interval := 0
				if res.Msg.R.Num != nil {
					totalNum = int(*res.Msg.R.Num)
				}
				if res.Msg.R.Interval != nil {
					interval = int(*res.Msg.R.Interval)
				}
				c.kTable.BatchCommand(ktable.PutPeer{ID: p.ID(), Addr: p.Addr(), Options: []ktable.PeerOption{
					ktable.PeerResponded(),
					ktable.PeerBep51Support(true),
					ktable.PeerSampleInfoHashesRes(
						len(hashesToStage),
						totalNum,
						time.Now().Add(time.Duration(interval)*time.Second),
					),
				}})
				if len(hashesToStage) > 0 {
					c.staging.stage(hashesToStage...)
				}
			})()
		}
	}
}
