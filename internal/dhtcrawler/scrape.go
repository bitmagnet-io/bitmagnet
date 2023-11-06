package dhtcrawler

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/bloom"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"time"
)

func (c *crawler) runScrape(ctx context.Context) {
	_ = c.scrape.Run(ctx, func(req nodeHasPeersForHash) {
		pfh, pfhErr := c.requestScrape(ctx, req)
		if pfhErr != nil {
			c.logger.Debugw("error getting peers scrape", "infoHash", req.infoHash.String(), "err", pfhErr)
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
	res, err := c.server.Query(ctx, req.node, dht.QGetPeers, dht.MsgArgs{
		ID:       c.kTable.Origin(),
		InfoHash: req.infoHash,
		Scrape:   1,
	})
	if err != nil {
		c.kTable.BatchCommand(ktable.DropAddr{Addr: req.node.Addr(), Reason: fmt.Errorf("failed to get peers from p: %w", err)})
		return infoHashWithScrape{}, err
	} else {
		c.kTable.BatchCommand(ktable.PutNode{ID: res.Msg.R.ID, Addr: req.node, Options: []ktable.NodeOption{ktable.NodeResponded()}})
	}
	if len(res.Msg.R.Nodes) > 0 {
		cancelCtx, cancel := context.WithTimeout(ctx, time.Second)
		for _, n := range res.Msg.R.Nodes {
			select {
			case <-cancelCtx.Done():
				break
			case c.discoveredNodes.In() <- ktable.NewNode(n.ID, n.Addr.ToAddrPort()):
				continue
			}
		}
		cancel()
	}
	if res.Msg.R.BFsd == nil || res.Msg.R.BFpe == nil {
		return infoHashWithScrape{}, errors.New("no scrape returned")
	}
	bfsd := bloom.FromScrape(*res.Msg.R.BFsd)
	bfpe := bloom.FromScrape(*res.Msg.R.BFpe)
	return infoHashWithScrape{
		nodeHasPeersForHash: req,
		bfsd:                bfsd,
		bfpe:                bfpe,
	}, nil
}
