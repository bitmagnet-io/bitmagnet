package dht_crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/channel"
	workers_metrics "github.com/bitmagnet-io/bitmagnet/internal/workers/metrics"
)

func newGetPeersWorker(
	cl client.Client,
	kTable ktable.Table,
	metaInfoRequestAdder channel.Adder[infoHashWithPeers],
	discoveredNodesAdder channel.Adder[ktable.Node],
	size atomic.Reader[GetPeersConcurrency],
	metrics *metrics.Component,
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
		channel.WithAtomicSize[nodeHasPeersForHash](
			atomic.MapIntish[GetPeersConcurrency, int](size),
		),
		channel.WithQuickShutdown[nodeHasPeersForHash](),
		channel.WithMetricsAdapter[nodeHasPeersForHash](
			workers_metrics.MustNew(metrics.MustSub(dht.QGetPeers)),
		),
	)
}
