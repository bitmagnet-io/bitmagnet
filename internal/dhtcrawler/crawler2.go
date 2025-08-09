package dhtcrawler

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/indexer"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
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
	metrics *metrics.Component,
) Runner {
	return runner.Runner(
		func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
			soughtNodeID := &concurrency.AtomicValue[protocol.ID]{}
			soughtNodeID.Set(protocol.RandomNodeID())

			enqueueProcessTorrents := newEnqueueProcessTorrentWorker(
				queueJobProvider,
				persisterAdder,
			)

			scrape := newScrapeWorker(
				client,
				kTable,
				persisterAdder,
				enqueueProcessTorrents,
				100,
				metrics,
			)

			persistTorrents := newPersistTorrentsWorker(
				classifier,
				persisterAdder,
				enqueueProcessTorrents,
				scrape,
				10,
				config.SavePieces,
				int(config.SaveFilesThreshold),
			)

			setupComplete := make(chan struct{})
			defer close(setupComplete)

			var discoveredNodesAdders []channel.Adder[ktable.Node]

			discoveredNodes := newDiscoveredNodesWorker(
				func() []channel.Adder[ktable.Node] {
					<-setupComplete

					return discoveredNodesAdders
				},
				time.Second*5,
				100,
			)

			bootstrap := newBootstrapper(
				config.BootstrapNodes,
				discoveredNodes,
			)

			ping := newPingWorker(
				client,
				kTable,
				100,
				time.Minute*15,
				metrics,
			)

			findNodes := newFindNodesWorker(
				client,
				kTable,
				discoveredNodes,
				soughtNodeID,
				1000,
				metrics,
				bootstrap,
			)

			requestMetaInfo := newRequestMetaInfoWorker(
				banningChecker,
				blockerBlocker,
				metainfoRequester,
				persistTorrents,
				1000,
				metrics,
			)

			getPeers := newGetPeersWorker(
				client,
				kTable,
				requestMetaInfo,
				discoveredNodes,
				1000,
				metrics,
			)

			infoHashTriage := newInfoHashTriageWorker(
				daoProvider,
				blockerBlocker,
				scrape,
				config.RescrapeThreshold,
				config.SaveFilesThreshold,
				getPeers,
				1000,
				time.Second*5,
				metrics,
			)

			sampleInfoHashes := newSampleInfoHashesWorker(
				client,
				kTable,
				infoHashTriage,
				discoveredNodes,
				soughtNodeID,
				1000,
				metrics,
				bootstrap,
			)

			discoveredNodesAdders = append(discoveredNodesAdders,
				sampleInfoHashes,
				findNodes,
				ping,
			)

			return concat.Runners(
				enqueueProcessTorrents,
				persistTorrents,
				discoveredNodes,
				ping,
				findNodes,
				newGetNodesForFindNodeWorker(kTable, findNodes),
				sampleInfoHashes,
				newNodesForSampleInfoHashesWorker(
					kTable,
					time.Second,
					sampleInfoHashes,
				),
				infoHashTriage,
				getPeers,
				requestMetaInfo,
				scrape,
				newOldNodesWorker(
					kTable,
					time.Second*10,
					time.Minute*15,
					ping,
				),
				periodic.New(
					config.ReseedBootstrapNodesInterval,
					bootstrap,
					periodic.WithInitialInterval(0),
					periodic.WithQuickShutdown(),
				),
				periodic.New(
					time.Minute,
					func(context.Context) error {
						soughtNodeID.Set(protocol.RandomNodeID())

						return nil
					},
				),
			)(ctx, cancel)
		})
}
