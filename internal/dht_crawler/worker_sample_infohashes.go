package dht_crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
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
	soughtNodeID *atomic.Value[protocol.ID],
	size atomic.Reader[SampleInfoHashesConcurrency],
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
				timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
				defer cancel()

				_ = discoveredNodesAdder.Add(timeoutCtx, slice.Map(res.Nodes, func(node client.NodeInfo) ktable.Node {
					return ktable.NewNode(node.ID, node.Addr)
				})...)
			}

			return nil
		},
		channel.WithAtomicSize[ktable.Node](
			atomic.MapIntish[SampleInfoHashesConcurrency, int](size),
		),
		channel.WithQuickShutdown[ktable.Node](),
		channel.WithMetricsAdapter[ktable.Node](
			workers_metrics.MustNew(metrics.MustSub(dht.QSampleInfohashes)),
		),
		// If this worker becomes idle we should bootstrap again:
		channel.WithOnIdle[ktable.Node](bootstrap),
	)
}
