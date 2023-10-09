package crawler

import (
	"context"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/bitmagnet-io/bitmagnet/internal/dht/staging"
)

// awaitInfoHashes waits for info hashes to be staged and forwards them to handleStagingRequest
func (c *crawler) awaitInfoHashes(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case req, ok := <-c.staging.Requested():
			if !ok {
				return
			}
			go c.handleStagingRequest(ctx, req)
		}
	}
}

func (c *crawler) handleStagingRequest(ctx context.Context, req staging.Request) {
	defer c.staging.Reject(req.InfoHash)
	pfh := newPeersForHash()
	pfhReq := peersForHashRequest{
		infoHash: req.InfoHash,
		peer:     req.Peer,
	}
	res := staging.Response{
		InfoHash: req.InfoHash,
	}
	defer func() {
		res.Scrape = staging.ResponseScrape{
			Bfsd:    pfh.bfsd,
			Bfpe:    pfh.bfpe,
			Scraped: pfh.scraped,
		}
		c.staging.Respond(ctx, res)
	}()
	if req.NeedMetaInfo {
		// need to get_peers with both scrape 0 and 1 because scrape responses usually omit peers values
		if err := c.requestPeersForHash(ctx, pfhReq, pfh); err != nil {
			c.logger.Debugw("error getting peers", "infoHash", req.InfoHash.String(), "err", err)
			return
		}
	}
	pfhReq.scrape = true
	if err := c.requestPeersForHash(ctx, pfhReq, pfh); err != nil {
		c.logger.Debugw("error getting peers with scrape", "infoHash", req.InfoHash.String(), "err", err)
		if c.discardUnscrapableTorrents || !req.NeedMetaInfo {
			return
		}
	}
	discoveredPeers := make([]krpc.NodeAddr, 0, len(pfh.peers))
	for _, p := range pfh.peers {
		discoveredPeers = append(discoveredPeers, p)
	}
	// Now we just need the torrent meta info; unfortunately we'll fail to get meta info for the vast majority of
	// discovered info hashes at this point, and the hashes will be discarded.
	// Maybe we could do something with some of these before discarding them, like requesting meta info from a caching service
	// such as iTorrents.org, though we get way too many to request them all...
	if req.NeedMetaInfo && res.MetaInfo.PieceLength == 0 {
		for _, p := range discoveredPeers {
			metaInfo, metaInfoErr := c.metainfoRequester.Request(ctx, req.InfoHash, p)
			if metaInfoErr == nil {
				res.MetaInfo = metaInfo
				break
			}
		}
		if res.MetaInfo.PieceLength == 0 {
			c.logger.Debugw("failed to get meta info from any peers", "infoHash", req.InfoHash.String())
		}
	}
	go c.routingTable.ReceivePeers(discoveredPeers...)
}
