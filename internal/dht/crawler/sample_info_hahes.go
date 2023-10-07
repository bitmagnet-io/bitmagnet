package crawler

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/bitmagnet-io/bitmagnet/internal/dht/routingtable"
	"github.com/bitmagnet-io/bitmagnet/internal/dht/staging"
	"time"
)

func (c *crawler) crawlPeersForInfoHashes(ctx context.Context) {
	fullLogged := false
	for {
		if c.staging.Count() < c.maxStagingSize {
			fullLogged = false
			go (func() {
				err := c.routingTable.TryEachNode(ctx, c.sampleInfoHashesFromLockedPeer)
				if err != nil {
					c.logger.Debugw("error crawling peers for info hashes", "err", err)
				}
			})()
		} else if !fullLogged {
			c.logger.Debug("staging is full, not crawling peers for info hashes")
			fullLogged = true
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.sampleInfoHashesInterval):
			continue
		}
	}
}

func (c *crawler) sampleInfoHashesFromLockedPeer(ctx context.Context, peer routingtable.PeerInfo) error {
	t := [20]byte{}
	if _, randErr := rand.Read(t[:]); randErr != nil {
		return fmt.Errorf("could not generate random bytes: %w", randErr)
	}
	res, resErr := c.dhtServer.Query(ctx, peer.Addr(), "sample_infohashes", krpc.MsgArgs{
		ID:     c.peerID,
		Target: t,
	})
	if resErr != nil {
		return resErr
	}
	if res.Msg.R == nil {
		return fmt.Errorf("sample_infohashes nil ret: %v", res.Msg)
	}
	c.routingTable.ReceiveNodeInfo(res.Msg.R.Nodes...)
	if res.Msg.R.Samples != nil {
		hashesToStage := make([]staging.InfoHashWithPeer, 0, len(*res.Msg.R.Samples))
		for _, s := range *res.Msg.R.Samples {
			hashesToStage = append(hashesToStage, staging.InfoHashWithPeer{
				InfoHash: s,
				Peer:     peer.Addr(),
			})
		}
		c.staging.Stage(hashesToStage...)
		c.logger.Debugw("staged hashes", "nHashes", len(hashesToStage))
	}
	return nil
}
