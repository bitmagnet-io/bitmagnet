package dhtcrawler

import (
	"context"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/prometheus/client_golang/prometheus"
)

func (c *crawler) runScrape(ctx context.Context) {
	_ = c.scrape.Run(ctx, func(req nodeHasPeersForHash) {
		pfh, pfhErr := c.requestScrape(ctx, req)
		if pfhErr != nil {
			return
		}
		select {
		case <-ctx.Done():
		case c.persistSources.In() <- infoHashWithScrape{
			nodeHasPeersForHash: req,
			bfsd:                pfh.bfsd,
			bfpe:                pfh.bfpe,
		}:
		}
	})
}

// requestScrape requests a scrape from a node to find seeders/leechers for a given info hash;
// see https://www.bittorrent.org/beps/bep_0033.html
func (c *crawler) requestScrape(
	ctx context.Context,
	req nodeHasPeersForHash,
) (infoHashWithScrape, error) {
	res, err := c.client.GetPeersScrape(ctx, req.node, req.infoHash)
	if err != nil {
		c.kTable.BatchCommand(ktable.DropAddr{
			Addr:   req.node.Addr(),
			Reason: fmt.Errorf("failed to get peers from p: %w", err),
		})

		return infoHashWithScrape{}, err
	}

	c.kTable.BatchCommand(ktable.PutNode{
		ID:      res.ID,
		Addr:    req.node,
		Options: []ktable.NodeOption{ktable.NodeResponded()},
	})

	c.scrapeNodeCount.Observe(float64(len(res.Nodes)))

	if len(res.Nodes) > 0 {
		cancelCtx, cancel := context.WithTimeout(ctx, time.Second)

		processed := 0

	nodes:
		for _, n := range res.Nodes {
			select {
			case <-cancelCtx.Done():
				break nodes
			case c.discoveredNodes.In() <- ktable.NewNode(n.ID, n.Addr):
				processed++
			}
		}

		c.scrapeNodeTotal.With(prometheus.Labels{"result": "discovered_nodes"}).Add(float64(processed))
		c.scrapeNodeTotal.With(prometheus.Labels{"result": "skipped"}).Add(float64(len(res.Nodes) - processed))

		cancel()
	}

	return infoHashWithScrape{
		nodeHasPeersForHash: req,
		bfsd:                res.BfSeeders,
		bfpe:                res.BfPeers,
	}, nil
}
