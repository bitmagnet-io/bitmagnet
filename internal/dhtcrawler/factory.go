package dhtcrawler

import (
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

type Runner runner.Provider

// func New(
// 	config Config,
// 	kTable ktable.Table,
// 	client client.Client,
// 	metainfoRequester metainforequester.Requester,
// 	banningChecker banning.Checker,
// 	daoProvider database.DaoTransactionProvider,
// 	queueJobProvider queue.JobProvider[processor.MessageParams],
// 	classifier classifier.Runner,
// 	persister persister.Adder,
// 	blocker blocking.Blocker,
// 	discoveredNodes concurrency.BatchingChannel[ktable.Node],
// 	logger *zap.SugaredLogger,
// ) (Runner, prometheus.Collector) {
// 	logger = logger.Named(Namespace)

// 	// persistedTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
// 	// 	Namespace: "bitmagnet",
// 	// 	Subsystem: Namespace,
// 	// 	Name:      "persisted_total",
// 	// 	Help:      "A counter of persisted database entities.",
// 	// }, []string{"entity"})

// 	isActive := &concurrency.AtomicValue[bool]{}

// 	run := func(ctx context.Context, cancel context.CancelCauseFunc) (shutdowner runner.Shutdowner, err error) {
// 		isActive.Set(true)

// 		defer func() {
// 			if err != nil {
// 				isActive.Set(false)
// 			}
// 		}()

// 		scalingFactor := int(config.ScalingFactor)

// 		ept := newEnqueueProcessTorrentWorker(
// 			queueJobProvider,
// 			persister,
// 		)

// 		pinger := newPingWorker(client, kTable, scalingFactor, time.Minute*15)

// 		scraper := newScrapeWorker(client, kTable, persister, ept, 10*scalingFactor)

// 		cr := &crawler{
// 			kTable:                       kTable,
// 			client:                       client,
// 			metainfoRequester:            metainfoRequester,
// 			banningChecker:               banningChecker,
// 			bootstrapNodes:               config.BootstrapNodes,
// 			reseedBootstrapNodesInterval: time.Minute * 10,
// 			getOldestNodesInterval:       time.Second * 10,
// 			// oldPeerThreshold:             time.Minute * 15,
// 			discoveredNodes: discoveredNodes,
// 			// nodesForPing: concurrency.NewBufferedConcurrentChannel[ktable.Node](
// 			// 	scalingFactor, scalingFactor),
// 			ping: pinger,
// 			nodesForFindNode: concurrency.NewBufferedConcurrentChannel[ktable.Node](
// 				10*scalingFactor, 10*scalingFactor),
// 			nodesForSampleInfoHashes: concurrency.NewBufferedConcurrentChannel[ktable.Node](
// 				10*scalingFactor,
// 				10*scalingFactor,
// 			),
// 			infoHashTriage: concurrency.NewBatchingChannel[nodesHavePeersForHash](
// 				10*scalingFactor, 1000, 20*time.Second),
// 			getPeers: concurrency.NewBufferedConcurrentChannel[nodesHavePeersForHash](
// 				10*scalingFactor, 20*scalingFactor),
// 			// scrape: concurrency.NewBufferedConcurrentChannel[nodeHasPeersForHash](
// 			// 	10*scalingFactor, 20*scalingFactor),
// 			requestMetaInfo: concurrency.NewBufferedConcurrentChannel[infoHashWithPeers](
// 				10*scalingFactor,
// 				40*scalingFactor,
// 			),
// 			persistTorrents: newPersistTorrentsWorker(
// 				10,
// 				classifier,
// 				persister,
// 				ept,
// 				scraper,
// 				config.SavePieces,
// 				int(config.SaveFilesThreshold),
// 				logger,
// 			),
// 			// classifyTorrents: concurrency.NewBufferedConcurrentChannel[infoHashWithMetaInfo](
// 			// 	10*scalingFactor,
// 			// 	10*scalingFactor,
// 			// ),
// 			enqueueProcessTorrents: ept,
// 			scrape:                 scraper,
// 			persistSources: concurrency.NewBatchingChannel[infoHashWithScrape](
// 				1000,
// 				1000,
// 				time.Minute,
// 			),
// 			classifier:         classifier,
// 			persister:          persister,
// 			saveFilesThreshold: config.SaveFilesThreshold,
// 			savePieces:         config.SavePieces,
// 			rescrapeThreshold:  config.RescrapeThreshold,
// 			daoProvider:        daoProvider,
// 			ignoreHashes: &ignoreHashes{
// 				bloom: boom.NewStableBloomFilter(10_000_000, 2, 0.001),
// 			},
// 			blockingManager:  blocker,
// 			soughtNodeID:     &concurrency.AtomicValue[protocol.ID]{},
// 			queueJobProvider: queueJobProvider,
// 			persistedTotal:   persistedTotal,
// 			logger:           logger,
// 		}

// 		shutdowner, err = cr.Runner(ctx, cancel)

// 		doShutdown := shutdowner
// 		shutdowner = func(ctx context.Context) error {
// 			isActive.Set(false)

// 			return doShutdown(ctx)
// 		}

// 		return
// 	}

// 	return run, persistedTotal
// }
