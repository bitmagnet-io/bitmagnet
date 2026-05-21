package dhtcrawler

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/prometheus/client_golang/prometheus"
)

func (c *crawler) runGetPeers(ctx context.Context) {
	_ = c.getPeers.Run(ctx, func(req nodeHasPeersForHash) {
		pfh, pfhErr := c.requestPeersForHash(ctx, req)
		if pfhErr != nil {
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
	res, err := c.client.GetPeers(ctx, req.node, req.infoHash)
	if err != nil {
		c.kTable.BatchCommand(ktable.DropAddr{
			Addr:   req.node.Addr(),
			Reason: fmt.Errorf("failed to get peers: %w", err),
		})

		return infoHashWithPeers{}, err
	}

	c.kTable.BatchCommand(ktable.PutNode{
		ID:      res.ID,
		Addr:    req.node,
		Options: []ktable.NodeOption{ktable.NodeResponded()},
	})

	c.getPeersPeerCount.Observe(float64(len(res.Values)))
	c.getPeersNodeCount.Observe(float64(len(res.Nodes)))

	if len(res.Nodes) > 0 {
		// block the channel for up to a second in an attempt to add the nodes to the discoveredNodes channel
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

		c.getPeersNodeTotal.With(prometheus.Labels{"result": "discovered_nodes"}).Add(float64(processed))
		c.getPeersNodeTotal.With(prometheus.Labels{"result": "skipped"}).Add(float64(len(res.Nodes) - processed))

		cancel()
	}

	if len(res.Values) < 1 {
		return infoHashWithPeers{}, errors.New("no peers found")
	}

	return infoHashWithPeers{
		nodeHasPeersForHash: req,
		peers:               res.Values,
	}, nil
}
