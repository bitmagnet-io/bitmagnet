package dhtcrawler

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"net/netip"
	"time"
)

func (c *crawler) reseedBootstrapNodes(ctx context.Context) {
	interval := time.Duration(0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			for _, addr := range c.bootstrapNodes {
				//if c.kTable.HasAddr(addr.Addr()) {
				//	continue
				//}
				_ = c.peersForFindNode.In(ctx, peer{addr: addr})
			}
		}
		interval = c.reseedBootstrapNodesInterval
	}
}

func (c *crawler) awaitPeersForPing(ctx context.Context) {
	args := dht.MsgArgs{
		ID: c.kTable.Origin(),
	}
	for {
		select {
		case <-ctx.Done():
			return
		case p := <-c.peersForPing.Out():
			go (func() {
				if p.Dropped() || p.Time().After(time.Now().Add(-c.oldPeerThreshold)) {
					// Either the p was already dropped or it succeeded after being added to the channel.
					// In either case we can continue.
					return
				}
				res, err := c.server.Query(ctx, p.Addr(), dht.QPing, args)
				var peerID protocol.ID
				if err == nil {
					peerID = res.Msg.R.ID
					if !p.ID().IsZero() && p.ID() != peerID {
						peerID = p.ID()
						err = errors.New("p responded with a mismatching ID")
					}
				}
				if err != nil {
					c.kTable.BatchCommand(ktable.DropPeer{ID: peerID, Reason: fmt.Errorf("failed to respond to ping: %w", err)})
				} else {
					c.kTable.BatchCommand(ktable.PutPeer{ID: peerID, Addr: p.Addr(), Options: []ktable.PeerOption{ktable.PeerResponded()}})
				}
			})()
		}
	}
}

func (c *crawler) getOldPeers(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.getOldestPeersInterval):
			for _, p := range c.kTable.GetOldestPeers(time.Now().Add(-c.oldPeerThreshold), 0) {
				_ = c.peersForPing.In(ctx, p)
			}
		}
	}
}

var _ ktable.Peer = peer{}

type peer struct {
	id   protocol.ID
	addr netip.AddrPort
}

func (p peer) ID() protocol.ID {
	return p.id
}

func (p peer) Addr() netip.AddrPort {
	return p.addr
}

func (p peer) Time() time.Time {
	return time.Time{}
}

func (p peer) Dropped() bool {
	return false
}

func (p peer) IsSampleInfoHashesCandidate() bool {
	return true
}
