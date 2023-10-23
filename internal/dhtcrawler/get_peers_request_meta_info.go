package dhtcrawler

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/bloom"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"net/netip"
)

// awaitInfoHashes waits for info hashes to be staged and forwards them to handleStagingRequest
func (c *crawler) awaitInfoHashes(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-c.staging.requested:
			go c.handleStagingRequest(ctx, req)
		}
	}
}

func (c *crawler) handleStagingRequest(ctx context.Context, req stagingRequest) {
	pfh := newPeersForHash()
	pfhReq := peersForHashRequest{
		infoHash: req.infoHash,
		peer:     req.peer,
	}
	res := stagingResponse{
		infoHash: req.infoHash,
	}
	defer func() {
		res.scrape = stagingResponseScrape{
			bfsd:    pfh.bfsd,
			bfpe:    pfh.bfpe,
			scraped: pfh.scraped,
		}
		c.staging.responseHolding <- res
	}()
	if req.needMetaInfo {
		// need to get_peers with both scrape 0 and 1 because scrape responses usually omit peers values
		if err := c.requestPeersForHash(ctx, pfhReq, pfh); err != nil {
			c.logger.Debugw("error getting peers", "infoHash", req.infoHash.String(), "err", err)
			return
		}
	}
	pfhReq.scrape = true
	if err := c.requestPeersForHash(ctx, pfhReq, pfh); err != nil {
		c.logger.Debugw("error getting peers with scrape", "infoHash", req.infoHash.String(), "err", err)
		if !req.needMetaInfo {
			return
		}
	}
	var hashPeers []ktable.HashPeer
	if req.needMetaInfo && res.metaInfo.PieceLength == 0 {
		var errs []error
		for _, p := range pfh.peers {
			var metaInfoErr error
			if res.metaInfo.PieceLength == 0 {
				if metaInfoRes, err := c.metainfoRequester.Request(ctx, req.infoHash, p); err == nil {
					res.metaInfo = metaInfoRes.Info
				} else {
					metaInfoErr = err
					errs = append(errs, err)
				}
			}
			if metaInfoErr == nil {
				hashPeers = append(hashPeers, ktable.HashPeer{
					Addr: p,
				})
				_ = c.discoveredPeers.TryIn(peer{
					addr: p,
				})
			}
		}
		if len(hashPeers) > 0 {
			c.kTable.BatchCommand(
				ktable.PutHash{ID: req.infoHash, Peers: hashPeers},
			)
		}
		if res.metaInfo.PieceLength == 0 {
			c.logger.Debugw("failed to get meta info from any peers", "infoHash", req.infoHash.String(), "err", errors.Join(errs...))
		}
	}
}

func (c *crawler) requestPeersForHash(
	ctx context.Context,
	req peersForHashRequest,
	pfh *peersForHash,
) error {
	scrape := 0
	if req.scrape {
		scrape = 1
	}
	res, err := c.server.Query(ctx, req.peer, dht.QGetPeers, dht.MsgArgs{
		ID:       c.kTable.Origin(),
		InfoHash: req.infoHash,
		Scrape:   scrape,
	})
	if err != nil {
		c.kTable.BatchCommand(ktable.DropAddr{Addr: req.peer.Addr(), Reason: fmt.Errorf("failed to get peers from p: %w", err)})
		return err
	} else {
		c.kTable.BatchCommand(ktable.PutPeer{ID: res.Msg.R.ID, Addr: req.peer, Options: []ktable.PeerOption{ktable.PeerResponded()}})
	}
	if res.Msg.R.BFsd != nil {
		thisBfsd := bloom.FromScrape(*res.Msg.R.BFsd)
		if bfsdErr := pfh.bfsd.Merge(&thisBfsd); bfsdErr != nil {
			return bfsdErr
		}
	}
	if res.Msg.R.BFpe != nil {
		thisBfpe := bloom.FromScrape(*res.Msg.R.BFpe)
		if bfpeErr := pfh.bfpe.Merge(&thisBfpe); bfpeErr != nil {
			return bfpeErr
		}
	}
	for _, p := range res.Msg.R.Values {
		if p.Port == 0 {
			continue
		}
		if _, ok := pfh.peers[p.String()]; !ok {
			pfh.peers[p.String()] = p.ToAddrPort()
		}
	}
	if req.scrape {
		pfh.scraped = true
	}
	return nil
}

type peersForHashRequest struct {
	infoHash protocol.ID
	peer     netip.AddrPort
	scrape   bool
}

type peersForHash struct {
	peers   map[string]netip.AddrPort
	bfsd    bloom.Filter
	bfpe    bloom.Filter
	scraped bool
}

func newPeersForHash() *peersForHash {
	return &peersForHash{
		peers: make(map[string]netip.AddrPort),
		bfsd:  bloom.New(),
		bfpe:  bloom.New(),
	}
}
