package crawler

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/bitmagnet-io/bitmagnet/internal/bloom"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type peersForHashRequest struct {
	infoHash model.Hash20
	peer     krpc.NodeAddr
	scrape   bool
}

type peersForHash struct {
	peers   map[string]krpc.NodeAddr
	bfsd    bloom.Filter
	bfpe    bloom.Filter
	scraped bool
}

func newPeersForHash() *peersForHash {
	return &peersForHash{
		peers: make(map[string]krpc.NodeAddr),
		bfsd:  bloom.New(),
		bfpe:  bloom.New(),
	}
}

func (c *crawler) requestPeersForHash(
	ctx context.Context,
	req peersForHashRequest,
	pfh *peersForHash,
) error {
	return c.routingTable.WithPeer(ctx, req.peer, func(ctx context.Context) error {
		t := [20]byte{}
		if _, randErr := rand.Read(t[:]); randErr != nil {
			return fmt.Errorf("could not generate random bytes: %w", randErr)
		}
		scrape := 0
		if req.scrape {
			scrape = 1
		}
		res, err := c.dhtServer.Query(ctx, req.peer, "get_peers", krpc.MsgArgs{
			ID:       c.peerID,
			InfoHash: krpc.ID(req.infoHash),
			Target:   t,
			Scrape:   scrape,
		})
		if err != nil {
			return err
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
		for _, peer := range res.Msg.R.Values {
			if peer.Port == 0 {
				continue
			}
			if _, ok := pfh.peers[peer.String()]; !ok {
				pfh.peers[peer.String()] = peer
			}
		}
		if req.scrape {
			pfh.scraped = true
		}
		return nil
	})
}
