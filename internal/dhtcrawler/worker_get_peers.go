package dhtcrawler

import (
	"context"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/channel"
)

func newGetPeersWorker(
	cl client.Client,
	kTable ktable.Table,
	metaInfoRequestAdder channel.Adder[infoHashWithPeers],
	discoveredNodesAdder channel.Adder[ktable.Node],
	size int,
) channel.Worker[nodeHasPeersForHash] {
	return channel.NewWorker(
		func(ctx context.Context, req nodeHasPeersForHash) error {
			res, err := cl.GetPeers(ctx, req.node, req.infoHash)
			if err != nil {
				err = fmt.Errorf("failed to get peers: %w", err)
				kTable.BatchCommand(ktable.DropAddr{
					Addr:   req.node.Addr(),
					Reason: err,
				})

				return nil
			}

			kTable.BatchCommand(ktable.PutNode{
				ID:      res.ID,
				Addr:    req.node,
				Options: []ktable.NodeOption{ktable.NodeResponded()},
			})

			if len(res.Nodes) > 0 {
				// block the channel for up to a second in an attempt to add the nodes to the discoveredNodes channel
				cancelCtx, cancel := context.WithTimeout(ctx, time.Second)

				discoveredNodesAdder.Add(cancelCtx, slice.Map(res.Nodes, func(info client.NodeInfo) ktable.Node {
					return ktable.NewNode(info.ID, info.Addr)
				})...)

				cancel()
			}

			if len(res.Values) == 0 {
				return nil
			}

			return metaInfoRequestAdder.Add(ctx, infoHashWithPeers{
				nodeHasPeersForHash: req,
				peers:               res.Values,
			})
		},
		channel.WithSize[nodeHasPeersForHash](size),
		channel.WithQuickShutdown[nodeHasPeersForHash](),
	)
}

// func (cr *crawler) runGetPeers(ctx context.Context) error {
// 	return cr.getPeers.Run(ctx, func(req nodeHasPeersForHash) {
// 		pfh, pfhErr := cr.requestPeersForHash(ctx, req)
// 		if pfhErr != nil {
// 			return
// 		}

// 		peers := make([]netip.AddrPort, 0, len(pfh.peers))
// 		hashPeers := make([]ktable.HashPeer, 0, len(pfh.peers))

// 		for _, p := range pfh.peers {
// 			peers = append(peers, p)
// 			hashPeers = append(hashPeers, ktable.HashPeer{
// 				Addr: p,
// 			})
// 		}

// 		cr.kTable.BatchCommand(
// 			ktable.PutHash{ID: req.infoHash, Peers: hashPeers},
// 		)
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case cr.requestMetaInfo.In() <- infoHashWithPeers{
// 			nodesHavePeersForHash: req,
// 			peers:                 peers,
// 		}:
// 			return
// 		}
// 	})
// }

// func requestPeersForHash(
// 	ctx context.Context,
// 	cl client.Client,
// 	kTable ktable.Table,
// 	discoveredNodesAdder channel.Adder[ktable.Node],
// 	req nodeHasPeersForHash,
// ) (infoHashWithPeers, error) {
// 	res, err := cl.GetPeers(ctx, req.node, req.infoHash)
// 	if err != nil {
// 		err = fmt.Errorf("failed to get peers: %w", err)
// 		kTable.BatchCommand(ktable.DropAddr{
// 			Addr:   req.node.Addr(),
// 			Reason: err,
// 		})

// 		return infoHashWithPeers{}, err
// 	}

// 	kTable.BatchCommand(ktable.PutNode{
// 		ID:      res.ID,
// 		Addr:    req.node,
// 		Options: []ktable.NodeOption{ktable.NodeResponded()},
// 	})

// 	if len(res.Nodes) > 0 {
// 		// block the channel for up to a second in an attempt to add the nodes to the discoveredNodes channel
// 		cancelCtx, cancel := context.WithTimeout(ctx, time.Second)

// 		discoveredNodesAdder.Add(cancelCtx, slice.Map(res.Nodes, func(info client.NodeInfo) ktable.Node {
// 			return ktable.NewNode(info.ID, info.Addr)
// 		})...)

// 		cancel()
// 	}

// 	if len(res.Values) == 0 {
// 		return infoHashWithPeers{}, errors.New("no peers found")
// 	}

// 	return infoHashWithPeers{
// 		nodeHasPeersForHash: req,
// 		peers:               res.Values,
// 	}, nil
// }
