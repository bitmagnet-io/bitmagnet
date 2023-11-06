package dhtcrawler

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"net/netip"
	"time"
)

func (c *crawler) runGetPeers(ctx context.Context) {
	_ = c.getPeers.Run(ctx, func(req nodeHasPeersForHash) {
		pfh, pfhErr := c.requestPeersForHash(ctx, req)
		if pfhErr != nil {
			c.logger.Debugw("error getting peers", "infoHash", req.infoHash.String(), "err", pfhErr)
			return
		}
		peers := make([]netip.AddrPort, 0, len(pfh.peers))
		hashPeers := make([]ktable.HashPeer, 0, len(pfh.peers))
		for _, p := range pfh.peers {
			peers = append(peers, p)
			hashPeers = append(hashPeers, ktable.HashPeer{
				Addr: p,
			})
		}
		c.kTable.BatchCommand(
			ktable.PutHash{ID: req.infoHash, Peers: hashPeers},
		)
		select {
		case <-ctx.Done():
			return
		case c.requestMetaInfo.In() <- infoHashWithPeers{
			nodeHasPeersForHash: req,
			peers:               peers,
		}:
			return
		}
	})
}

func (c *crawler) requestPeersForHash(
	ctx context.Context,
	req nodeHasPeersForHash,
) (infoHashWithPeers, error) {
	res, err := c.server.Query(ctx, req.node, dht.QGetPeers, dht.MsgArgs{
		ID:       c.kTable.Origin(),
		InfoHash: req.infoHash,
	})
	if err != nil {
		c.kTable.BatchCommand(ktable.DropAddr{Addr: req.node.Addr(), Reason: fmt.Errorf("failed to get peers: %w", err)})
		return infoHashWithPeers{}, err
	} else {
		c.kTable.BatchCommand(ktable.PutNode{ID: res.Msg.R.ID, Addr: req.node, Options: []ktable.NodeOption{ktable.NodeResponded()}})
	}
	if len(res.Msg.R.Nodes) > 0 {
		// block the channel for up to a second in an attempt to add the nodes to the discoveredNodes channel
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
	if len(res.Msg.R.Values) < 1 {
		return infoHashWithPeers{}, errors.New("no peers found")
	}
	addrs := make([]netip.AddrPort, 0, len(res.Msg.R.Values))
	for _, p := range res.Msg.R.Values {
		if p.Port == 0 {
			continue
		}
		addrs = append(addrs, p.ToAddrPort())
	}
	return infoHashWithPeers{
		nodeHasPeersForHash: req,
		peers:               addrs,
	}, nil
}
