package dhtcrawler

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/blocking"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/dhtcrawler/dhtcrawlerhealthcheck"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/banning"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainforequester"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"github.com/prometheus/client_golang/prometheus"
	boom "github.com/tylertreat/BoomFilters"
	"go.uber.org/zap"
)

func New(
	config Config,
	kTable ktable.Table,
	client client.Client,
	metainfoRequester metainforequester.Requester,
	banningChecker banning.Checker,
	daoProvider database.DaoTransactionProvider,
	blocker blocking.Blocker,
	discoveredNodes concurrency.BatchingChannel[ktable.Node],
	dhtServerLastResponses *concurrency.AtomicValue[server.LastResponses],
	logger *zap.SugaredLogger,
) (registry.Option, health.CheckerOption, prometheus.Collector) {
	logger = logger.Named(Namespace)

	persistedTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "bitmagnet",
		Subsystem: Namespace,
		Name:      "persisted_total",
		Help:      "A counter of persisted database entities.",
	}, []string{"entity"})

	isActive := &concurrency.AtomicValue[bool]{}

	run := func(ctx context.Context, cancel context.CancelCauseFunc) (shutdowner runner.Shutdowner, err error) {
		isActive.Set(true)

		defer func() {
			if err != nil {
				isActive.Set(false)
			}
		}()

		scalingFactor := int(config.ScalingFactor)

		cr := &crawler{
			kTable:                       kTable,
			client:                       client,
			metainfoRequester:            metainfoRequester,
			banningChecker:               banningChecker,
			bootstrapNodes:               config.BootstrapNodes,
			reseedBootstrapNodesInterval: time.Minute * 10,
			getOldestNodesInterval:       time.Second * 10,
			oldPeerThreshold:             time.Minute * 15,
			discoveredNodes:              discoveredNodes,
			nodesForPing: concurrency.NewBufferedConcurrentChannel[ktable.Node](
				scalingFactor, scalingFactor),
			nodesForFindNode: concurrency.NewBufferedConcurrentChannel[ktable.Node](
				10*scalingFactor, 10*scalingFactor),
			nodesForSampleInfoHashes: concurrency.NewBufferedConcurrentChannel[ktable.Node](
				10*scalingFactor,
				10*scalingFactor,
			),
			infoHashTriage: concurrency.NewBatchingChannel[nodeHasPeersForHash](
				10*scalingFactor, 1000, 20*time.Second),
			getPeers: concurrency.NewBufferedConcurrentChannel[nodeHasPeersForHash](
				10*scalingFactor, 20*scalingFactor),
			scrape: concurrency.NewBufferedConcurrentChannel[nodeHasPeersForHash](
				10*scalingFactor, 20*scalingFactor),
			requestMetaInfo: concurrency.NewBufferedConcurrentChannel[infoHashWithPeers](
				10*scalingFactor,
				40*scalingFactor,
			),
			persistTorrents: concurrency.NewBatchingChannel[infoHashWithMetaInfo](
				1000,
				1000,
				time.Minute,
			),
			persistSources: concurrency.NewBatchingChannel[infoHashWithScrape](
				1000,
				1000,
				time.Minute,
			),
			saveFilesThreshold: config.SaveFilesThreshold,
			savePieces:         config.SavePieces,
			rescrapeThreshold:  config.RescrapeThreshold,
			daoProvider:        daoProvider,
			ignoreHashes: &ignoreHashes{
				bloom: boom.NewStableBloomFilter(10_000_000, 2, 0.001),
			},
			blockingManager: blocker,
			soughtNodeID:    &concurrency.AtomicValue[protocol.ID]{},
			persistedTotal:  persistedTotal,
			logger:          logger,
		}

		shutdowner, err = cr.Runner(ctx, cancel)

		doShutdown := shutdowner
		shutdowner = func(ctx context.Context) error {
			isActive.Set(false)

			return doShutdown(ctx)
		}

		return
	}

	return registry.WithWorker(
			Namespace,
			run,
			worker.WithDependencies(server.Namespace, database.Namespace, blocking.Namespace),
		),
		dhtcrawlerhealthcheck.New(dhtServerLastResponses, isActive),
		persistedTotal
}
