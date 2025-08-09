package dhtcrawler

import (
	"context"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/batch"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/channel"
	workers_metrics "github.com/bitmagnet-io/bitmagnet/internal/workers/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/periodic"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

func newNodesForSampleInfoHashesWorker(
	kTable ktable.Table,
	interval time.Duration,
	sampleInfoHashesAdder channel.Adder[ktable.Node],
) runner.Provider {
	return periodic.New(
		interval,
		func(ctx context.Context) error {
			return sampleInfoHashesAdder.Add(ctx, kTable.GetNodesForSampleInfoHashes(60)...)
		},
	)
}

func newSampleInfoHashesWorker(
	cl client.Client,
	kTable ktable.Table,
	infoHashTriageAdder batch.Adder[nodeHasPeersForHash],
	discoveredNodesAdder channel.Adder[ktable.Node],
	soughtNodeID *concurrency.AtomicValue[protocol.ID],
	size int,
	metrics *metrics.Component,
	bootstrap func(context.Context) error,
) channel.Worker[ktable.Node] {
	return channel.NewWorker(
		func(ctx context.Context, node ktable.Node) error {
			if !node.IsSampleInfoHashesCandidate() {
				return nil
			}

			res, err := cl.SampleInfoHashes(ctx, node.Addr(), soughtNodeID.Get())
			if err != nil {
				kTable.BatchCommand(
					ktable.DropNode{ID: node.ID(), Reason: fmt.Errorf("sample_infohashes failed: %w", err)},
				)

				return nil
			}

			infoHashTriageAdder.Add(ctx, slice.Map(res.Samples, func(hash protocol.ID) nodeHasPeersForHash {
				return nodeHasPeersForHash{
					infoHash: hash,
					node:     node.Addr(),
				}
			})...)

			interval := res.Interval
			// most nodes request a 6 hour backoff time(!)
			// if we're still discovering info hashes from them then let's set a respectful interval instead
			if len(res.Samples) > 0 && interval > 300 {
				interval = 60
			}

			kTable.BatchCommand(ktable.PutNode{ID: node.ID(), Addr: node.Addr(), Options: []ktable.NodeOption{
				ktable.NodeResponded(),
				ktable.NodeBep51Support(true),
				ktable.NodeSampleInfoHashesRes(
					len(res.Samples),
					res.Num,
					time.Now().Add(time.Duration(interval)*time.Second),
				),
			}})

			if len(res.Nodes) > 0 {
				// block on the channel for up to a second trying to add sampled nodes to the discoveredNodes
				// channel
				// go func() {
				timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
				defer cancel()

				_ = discoveredNodesAdder.Add(timeoutCtx, slice.Map(res.Nodes, func(node client.NodeInfo) ktable.Node {
					return ktable.NewNode(node.ID, node.Addr)
				})...)
			}

			return nil
		},
		channel.WithSize[ktable.Node](size),
		channel.WithQuickShutdown[ktable.Node](),
		channel.WithMetricsAdapter[ktable.Node](
			workers_metrics.MustNew(metrics.MustSub(dht.QSampleInfohashes)),
		),
		// If this worker becomes idle we should bootstrap again:
		channel.WithOnIdle[ktable.Node](bootstrap),
	)
}

// func (cr *crawler) getNodesForSampleInfoHashes(ctx context.Context) error {
// 	for {
// 		peers := cr.kTable.GetNodesForSampleInfoHashes(60)
// 		for _, p := range peers {
// 			select {
// 			case <-ctx.Done():
// 				return ctx.Err()
// 			case cr.nodesForSampleInfoHashes.In() <- p:
// 				continue
// 			}
// 		}

// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		case <-time.After(time.Second):
// 		}
// 	}
// }

// func (cr *crawler) runSampleInfoHashes(ctx context.Context) error {
// 	return cr.nodesForSampleInfoHashes.Run(ctx, func(n ktable.Node) {
// 		if !n.IsSampleInfoHashesCandidate() {
// 			return
// 		}

// 		res, err := cr.client.SampleInfoHashes(ctx, n.Addr(), cr.soughtNodeID.Get())
// 		if err != nil {
// 			cr.kTable.BatchCommand(
// 				ktable.DropNode{ID: n.ID(), Reason: fmt.Errorf("sample_infohashes failed: %w", err)},
// 			)

// 			return
// 		}

// 		var discoveredHashes []nodeHasPeersForHash

// 		for _, s := range res.Samples {
// 			if !cr.ignoreHashes.testAndAdd(s) {
// 				discoveredHashes = append(discoveredHashes, nodeHasPeersForHash{
// 					infoHash: s,
// 					node:     n.Addr(),
// 				})
// 			}
// 		}

// 		for _, h := range discoveredHashes {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case cr.infoHashTriage.In() <- h:
// 				continue
// 			}
// 		}

// 		interval := res.Interval
// 		// most nodes request a 6 hour backoff time(!)
// 		// if we're still discovering info hashes from them then let's set a respectful interval instead
// 		if len(discoveredHashes) > 0 && interval > 300 {
// 			interval = 60
// 		}

// 		cr.kTable.BatchCommand(ktable.PutNode{ID: n.ID(), Addr: n.Addr(), Options: []ktable.NodeOption{
// 			ktable.NodeResponded(),
// 			ktable.NodeBep51Support(true),
// 			ktable.NodeSampleInfoHashesRes(
// 				len(discoveredHashes),
// 				res.Num,
// 				time.Now().Add(time.Duration(interval)*time.Second),
// 			),
// 		}})

// 		if len(res.Nodes) > 0 {
// 			// block on the channel for up to a second trying to add sampled nodes to the discoveredNodes
// 			// channel
// 			go func() {
// 				timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
// 				defer cancel()

// 				for _, n := range res.Nodes {
// 					select {
// 					case <-timeoutCtx.Done():
// 						return
// 					case cr.discoveredNodes.In() <- ktable.NewNode(n.ID, n.Addr):
// 						continue
// 					}
// 				}
// 			}()
// 		}
// 	})
// }
