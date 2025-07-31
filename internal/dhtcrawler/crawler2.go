package dhtcrawler

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/indexer"
	"github.com/bitmagnet-io/bitmagnet/internal/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/banning"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainforequester"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/channel"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/concat"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/periodic"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
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
	queueJobProvider queue.JobProvider[indexer.MessageParams],
	classifier classifier.Runner,
	persisterAdder persister.Adder,
	blockerBlocker blocker.Blocker,
	logger *zap.SugaredLogger,
) Runner {
	return runner.Runner(
		func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
			setupComplete := make(chan struct{})
			defer close(setupComplete)

			soughtNodeID := &concurrency.AtomicValue[protocol.ID]{}
			soughtNodeID.Set(protocol.RandomNodeID())

			enqueueProcessTorrents := newEnqueueProcessTorrentWorker(
				queueJobProvider,
				persisterAdder,
			)

			scraper := newScrapeWorker(
				client,
				kTable,
				persisterAdder,
				enqueueProcessTorrents,
				100,
			)

			persistTorrents := newPersistTorrentsWorker(
				classifier,
				persisterAdder,
				enqueueProcessTorrents,
				scraper,
				10,
				config.SavePieces,
				int(config.SaveFilesThreshold),
			)

			ignoreHashes := &ignoreHashes{
				bloom: boom.NewStableBloomFilter(10_000_000, 2, 0.001),
			}

			var discoveredNodesAdders []channel.Adder[ktable.Node]

			discoveredNodes := newDiscoveredNodesWorker(
				func() []channel.Adder[ktable.Node] {
					<-setupComplete

					return discoveredNodesAdders
				},
				time.Second*5,
				100,
			)

			ping := newPingWorker(
				client,
				kTable,
				100,
				time.Minute*15,
			)

			findNodes := newFindNodesWorker(
				client,
				kTable,
				discoveredNodes,
				soughtNodeID,
				1000,
			)

			requestMetaInfo := newRequestMetaInfoWorker(
				banningChecker,
				blockerBlocker,
				metainfoRequester,
				persistTorrents,
				1000,
			)

			getPeers := newGetPeersWorker(
				client,
				kTable,
				requestMetaInfo,
				discoveredNodes,
				1000,
			)

			infoHashTriage := newInfoHashTriageWorker(
				daoProvider,
				ignoreHashes,
				blockerBlocker,
				scraper,
				config.RescrapeThreshold,
				config.SaveFilesThreshold,
				getPeers,
				1000,
				time.Second*5,
			)

			sampleInfoHashes := newSampleInfoHashesWorker(
				client,
				kTable,
				infoHashTriage,
				discoveredNodes,
				soughtNodeID,
				1000,
			)

			discoveredNodesAdders = append(discoveredNodesAdders,
				sampleInfoHashes,
				findNodes,
				ping,
			)

			return concat.Runners(
				enqueueProcessTorrents.Runner(),
				persistTorrents.Runner(),
				discoveredNodes.Runner(),
				ping.Runner(),
				findNodes.Runner(),
				newGetNodesForFindNodeWorker(kTable, findNodes).Runner(),
				sampleInfoHashes.Runner(),
				newNodesForSampleInfoHashesWorker(
					kTable,
					time.Second,
					sampleInfoHashes,
				).Runner(),
				infoHashTriage.Runner(),
				getPeers.Runner(),
				requestMetaInfo.Runner(),
				scraper.Runner(),
				newBootstrapNodesWorker(
					config.ReseedBootstrapNodesInterval,
					config.BootstrapNodes,
					discoveredNodes,
					logger,
				).Runner(),
				newOldNodesWorker(
					kTable,
					time.Second*10,
					time.Minute*15,
					ping,
				).Runner(),
				periodic.New(
					time.Minute,
					func(context.Context) error {
						soughtNodeID.Set(protocol.RandomNodeID())

						return nil
					},
				).Runner(),
			).Runner()(ctx, cancel)
		})
}
