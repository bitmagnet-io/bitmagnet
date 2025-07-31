package dhtcrawler

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/batch"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/channel"
)

func newScrapeWorker(
	cl client.Client,
	kTable ktable.Table,
	persisterAdder persister.Adder,
	processorAdder batch.Adder[protocol.ID],
	size int,
) channel.Worker[nodeHasPeersForHash] {
	return channel.NewWorker(
		func(ctx context.Context, req nodeHasPeersForHash) error {
			res, err := cl.GetPeersScrape(ctx, req.node, req.infoHash)
			if err != nil {
				kTable.BatchCommand(ktable.DropAddr{
					Addr:   req.node.Addr(),
					Reason: fmt.Errorf("failed to get peers from p: %w", err),
				})

				return nil
			}

			kTable.BatchCommand(ktable.PutNode{
				ID:      res.ID,
				Addr:    req.node,
				Options: []ktable.NodeOption{ktable.NodeResponded()},
			})

			err = persisterAdder.Add(
				ctx,
				persister.InputTorrentsTorrentSources(
					model.TorrentsTorrentSource{
						Source:   "dht",
						InfoHash: req.infoHash,
						Seeders:  model.NewNullUint(uint(res.BfSeeders.ApproximatedSize())),
						Leechers: model.NewNullUint(uint(res.BfPeers.ApproximatedSize())),
					},
				),
			)

			if err != nil {
				return err
			}

			return processorAdder.Add(ctx, req.infoHash)
		},
		channel.WithSize[nodeHasPeersForHash](size),
		channel.WithQuickShutdown[nodeHasPeersForHash](),
	)
}

// requestScrape requests a scrape from a node to find seeders/leechers for a given info hash;
// see https://www.bittorrent.org/beps/bep_0033.html
// func requestScrape(
// 	ctx context.Context,
// 	cl client.Client,
// 	kTable ktable.Table,
// 	req nodeHasPeersForHash,
// ) (infoHashWithScrape, error) {
// 	res, err := cl.GetPeersScrape(ctx, req.node, req.infoHash)
// 	if err != nil {
// 		kTable.BatchCommand(ktable.DropAddr{
// 			Addr:   req.node.Addr(),
// 			Reason: fmt.Errorf("failed to get peers from p: %w", err),
// 		})

// 		return infoHashWithScrape{}, err
// 	}

// 	kTable.BatchCommand(ktable.PutNode{
// 		ID:      res.ID,
// 		Addr:    req.node,
// 		Options: []ktable.NodeOption{ktable.NodeResponded()},
// 	})

// 	return infoHashWithScrape{
// 		nodeHasPeersForHash: req,
// 		bfsd:                res.BfSeeders,
// 		bfpe:                res.BfPeers,
// 	}, nil
// }
