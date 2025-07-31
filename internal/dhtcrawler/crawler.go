package dhtcrawler

import (
	"net/netip"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/bloom"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
	boom "github.com/tylertreat/BoomFilters"
)

// const Namespace = "dht_crawler"

// //type Crawler interface {
// //	runner.Interface
// //	//IsActive() bool
// //}

// type crawler struct {
// 	kTable                       ktable.Table
// 	client                       client.Client
// 	metainfoRequester            metainforequester.Requester
// 	banningChecker               banning.Checker
// 	bootstrapNodes               []string
// 	reseedBootstrapNodesInterval time.Duration
// 	getOldestNodesInterval       time.Duration
// 	oldNodeThreshold             time.Duration
// 	discoveredNodes              concurrency.BatchingChannel[ktable.Node]
// 	// nodesForPing                 concurrency.BufferedConcurrentChannel[ktable.Node]
// 	nodesForFindNode         concurrency.BufferedConcurrentChannel[ktable.Node]
// 	nodesForSampleInfoHashes concurrency.BufferedConcurrentChannel[ktable.Node]
// 	infoHashTriage           concurrency.BatchingChannel[nodesHavePeersForHash]
// 	getPeers                 concurrency.BufferedConcurrentChannel[nodesHavePeersForHash]
// 	// scrape                       concurrency.BufferedConcurrentChannel[nodeHasPeersForHash]
// 	requestMetaInfo        concurrency.BufferedConcurrentChannel[infoHashWithPeers]
// 	persistTorrents        channel.Worker[infoHashWithMetaInfo]
// 	enqueueProcessTorrents batch.Worker[protocol.ID]
// 	ping                   channel.Worker[ktable.Node]
// 	scrape                 channel.Worker[nodesHavePeersForHash]
// 	persistSources         concurrency.BatchingChannel[infoHashWithScrape]
// 	rescrapeThreshold      time.Duration
// 	saveFilesThreshold     uint
// 	savePieces             bool
// 	daoProvider            database.DaoTransactionProvider
// 	// ignoreHashes is a thread-safe bloom filter that the crawler keeps in memory,
// 	// containing every hash it has already encountered.
// 	// This avoids multiple attempts to crawl the same hash, and takes a lot of load off the database query
// 	// that checks if a hash has already been indexed.
// 	ignoreHashes    *ignoreHashes
// 	blockingManager blocking.Blocker
// 	// soughtNodeID is a random node ID used as the target for find_node and sample_infohashes requests.
// 	// It is rotated every 10 seconds.
// 	soughtNodeID     *concurrency.AtomicValue[protocol.ID]
// 	persistedTotal   *prometheus.CounterVec
// 	logger           *zap.SugaredLogger
// 	queueJobProvider queue.JobProvider[processor.MessageParams]
// 	classifier       classifier.Runner
// 	persister        persister.Adder
// }

// func (cr *crawler) Runner(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
// 	return concat.Runners(
// 		cr.runPersistSources,
// 		cr.enqueueProcessTorrents.Runner,
// 		cr.persistTorrents.Runner,
// 		runner.SimpleRunner(cr.rotateSoughtNodeID),
// 		runner.SimpleRunner(cr.runDiscoveredNodes),
// 		cr.ping.Runner,
// 		runner.SimpleRunner(cr.runFindNode),
// 		runner.SimpleRunner(cr.getNodesForFindNode),
// 		runner.SimpleRunner(cr.runSampleInfoHashes),
// 		runner.SimpleRunner(cr.getNodesForSampleInfoHashes),
// 		runner.SimpleRunner(cr.runInfoHashTriage),
// 		runner.SimpleRunner(cr.runGetPeers),
// 		runner.SimpleRunner(cr.runRequestMetaInfo),
// 		cr.scrape.Runner,
// 		newBootstrapNodesWorker(cr.reseedBootstrapNodesInterval, cr.bootstrapNodes, cr.ping, cr.logger).Runner,
// 		newOldNodesWorker(cr.kTable, cr.getOldestNodesInterval, cr.oldNodeThreshold, cr.ping).Runner,
// 	)(ctx, cancel)
// }

type nodeHasPeersForHash struct {
	infoHash               protocol.ID
	node                   netip.AddrPort
	isVerifiedAbsentFromDB bool
}

type infoHashWithMetaInfo struct {
	nodeHasPeersForHash
	metaInfo metainfo.Info
}

type infoHashWithPeers struct {
	nodeHasPeersForHash
	peers []netip.AddrPort
}

type infoHashWithScrape struct {
	nodeHasPeersForHash
	bfsd bloom.Filter
	bfpe bloom.Filter
}

type ignoreHashes struct {
	mutex sync.Mutex
	bloom *boom.StableBloomFilter
}

func (i *ignoreHashes) testAndAdd(id protocol.ID) bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	return i.bloom.TestAndAdd(id[:])
}

// func (cr *crawler) rotateSoughtNodeID(ctx context.Context) error {
// 	interval := time.Duration(0)

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		case <-time.After(interval):
// 			cr.soughtNodeID.Set(protocol.RandomNodeID())
// 		}

// 		interval = time.Minute
// 	}
// }
