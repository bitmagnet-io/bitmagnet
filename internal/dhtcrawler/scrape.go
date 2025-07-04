package dhtcrawler

import (
	"context"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
)

func (cr *crawler) runScrape(ctx context.Context) error {
	return cr.scrape.Run(ctx, func(req nodeHasPeersForHash) {
		pfh, pfhErr := cr.requestScrape(ctx, req)
		if pfhErr != nil {
			return
		}
		select {
		case <-ctx.Done():
		case cr.persistSources.In() <- infoHashWithScrape{
			nodeHasPeersForHash: req,
			bfsd:                pfh.bfsd,
			bfpe:                pfh.bfpe,
		}:
		}
	})
}

// requestScrape requests a scrape from a node to find seeders/leechers for a given info hash;
// see https://www.bittorrent.org/beps/bep_0033.html
func (cr *crawler) requestScrape(
	ctx context.Context,
	req nodeHasPeersForHash,
) (infoHashWithScrape, error) {
	res, err := cr.client.GetPeersScrape(ctx, req.node, req.infoHash)
	if err != nil {
		cr.kTable.BatchCommand(ktable.DropAddr{
			Addr:   req.node.Addr(),
			Reason: fmt.Errorf("failed to get peers from p: %w", err),
		})

		return infoHashWithScrape{}, err
	}

	cr.kTable.BatchCommand(ktable.PutNode{
		ID:      res.ID,
		Addr:    req.node,
		Options: []ktable.NodeOption{ktable.NodeResponded()},
	})

	if len(res.Nodes) > 0 {
		cancelCtx, cancel := context.WithTimeout(ctx, time.Second)

		for _, n := range res.Nodes {
			select {
			case <-cancelCtx.Done():
				break
			case cr.discoveredNodes.In() <- ktable.NewNode(n.ID, n.Addr):
				continue
			}
		}

		cancel()
	}

	return infoHashWithScrape{
		nodeHasPeersForHash: req,
		bfsd:                res.BfSeeders,
		bfpe:                res.BfPeers,
	}, nil
}
