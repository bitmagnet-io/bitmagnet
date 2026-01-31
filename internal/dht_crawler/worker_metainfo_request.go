package dht_crawler

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/banning"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainforequester"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/channel"
	workers_metrics "github.com/bitmagnet-io/bitmagnet/internal/workers/metrics"
)

const MetricRequestMetaInfo = "request_meta_info"

func newRequestMetaInfoWorker(
	banningChecker banning.Checker,
	blockerBlocker blocker.Blocker,
	metainfoRequester metainforequester.Requester,
	persistTorrentsAdder channel.Adder[infoHashWithMetaInfo],
	size atomic.Reader[RequestMetaInfoConcurrency],
	metrics *metrics.Component,
) channel.Worker[infoHashWithPeers] {
	return channel.NewWorker(
		func(ctx context.Context, req infoHashWithPeers) error {
			var (
				response    metainforequester.Response
				gotResponse bool
				err         error
			)

			for _, peer := range req.peers {
				response, err = metainfoRequester.Request(ctx, req.infoHash, peer)
				if err == nil {
					gotResponse = true
					break
				}
			}

			if !gotResponse {
				return nil
			}

			if err = banningChecker.Check(response.Info); err != nil {
				return blockerBlocker.Block(ctx, []protocol.ID{req.infoHash}, false)
			}

			return persistTorrentsAdder.Add(ctx, infoHashWithMetaInfo{
				nodeHasPeersForHash: req.nodeHasPeersForHash,
				metaInfo:            response.Info,
			})
		},
		channel.WithAtomicSize[infoHashWithPeers](
			atomic.MapIntish[RequestMetaInfoConcurrency, int](size),
		),
		channel.WithQuickShutdown[infoHashWithPeers](),
		channel.WithMetricsAdapter[infoHashWithPeers](
			workers_metrics.MustNew(metrics.MustSub(MetricRequestMetaInfo)),
		),
	)
}
