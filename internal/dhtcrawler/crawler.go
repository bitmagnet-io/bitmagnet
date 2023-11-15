package dhtcrawler

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/bloom"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/worker"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/message"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/banning"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainforequester"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/publisher"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/netip"
	"sync"
	"time"
)

type Params struct {
	fx.In
	Config              Config
	KTable              ktable.Table
	Client              client.Client
	MetainfoRequester   metainforequester.Requester
	BanningChecker      banning.Checker `name:"metainfo_banning_checker"`
	Search              search.Search
	Dao                 *dao.Query
	ClassifierPublisher publisher.Publisher[message.ClassifyTorrentPayload]
	DiscoveredNodes     concurrency.BatchingChannel[ktable.Node] `name:"dht_discovered_nodes"`
	Logger              *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Worker worker.Worker `group:"workers"`
}

func New(params Params) Result {
	scalingFactor := int(params.Config.ScalingFactor)
	c := crawler{
		kTable:                       params.KTable,
		client:                       params.Client,
		metainfoRequester:            params.MetainfoRequester,
		banningChecker:               params.BanningChecker,
		bootstrapNodes:               params.Config.BootstrapNodes,
		reseedBootstrapNodesInterval: time.Minute * 10,
		getOldestNodesInterval:       time.Second * 10,
		oldPeerThreshold:             time.Minute * 15,
		discoveredNodes:              params.DiscoveredNodes,
		nodesForPing:                 concurrency.NewBufferedConcurrentChannel[ktable.Node](scalingFactor, scalingFactor),
		nodesForFindNode:             concurrency.NewBufferedConcurrentChannel[ktable.Node](10*scalingFactor, 10*scalingFactor),
		nodesForSampleInfoHashes:     concurrency.NewBufferedConcurrentChannel[ktable.Node](10*scalingFactor, 100*scalingFactor),
		infoHashTriage:               concurrency.NewBatchingChannel[nodeHasPeersForHash](100*scalingFactor, 100, 5*time.Second),
		getPeers:                     concurrency.NewBufferedConcurrentChannel[nodeHasPeersForHash](10*scalingFactor, 100*scalingFactor),
		scrape:                       concurrency.NewBufferedConcurrentChannel[nodeHasPeersForHash](10*scalingFactor, 100*scalingFactor),
		requestMetaInfo:              concurrency.NewBufferedConcurrentChannel[infoHashWithPeers](10*scalingFactor, 100*scalingFactor),
		persistTorrents: concurrency.NewBatchingChannel[infoHashWithMetaInfo](
			1000,
			100,
			10*time.Second,
		),
		persistSources: concurrency.NewBatchingChannel[infoHashWithScrape](
			1000,
			100,
			10*time.Second,
		),
		saveFilesThreshold:  params.Config.SaveFilesThreshold,
		savePieces:          params.Config.SavePieces,
		rescrapeThreshold:   params.Config.RescrapeThreshold,
		dao:                 params.Dao,
		classifierPublisher: params.ClassifierPublisher,
		ignoreHashes: &ignoreHashes{
			bloom: *bloom.NewWithEstimates(1000000, 0.0001),
		},
		soughtNodeID: &concurrency.AtomicValue[protocol.ID]{},
		stopped:      make(chan struct{}),
		logger:       params.Logger.Named("dht_crawler"),
	}
	c.soughtNodeID.Set(protocol.RandomNodeID())
	return Result{
		Worker: worker.NewWorker(
			"dht_crawler",
			fx.Hook{
				OnStart: func(context.Context) error {
					go c.start()
					return nil
				},
				OnStop: func(context.Context) error {
					close(c.stopped)
					return nil
				},
			},
		),
	}
}

type crawler struct {
	kTable                       ktable.Table
	client                       client.Client
	metainfoRequester            metainforequester.Requester
	banningChecker               banning.Checker
	bootstrapNodes               []string
	reseedBootstrapNodesInterval time.Duration
	getOldestNodesInterval       time.Duration
	oldPeerThreshold             time.Duration
	discoveredNodes              concurrency.BatchingChannel[ktable.Node]
	nodesForPing                 concurrency.BufferedConcurrentChannel[ktable.Node]
	nodesForFindNode             concurrency.BufferedConcurrentChannel[ktable.Node]
	nodesForSampleInfoHashes     concurrency.BufferedConcurrentChannel[ktable.Node]
	infoHashTriage               concurrency.BatchingChannel[nodeHasPeersForHash]
	getPeers                     concurrency.BufferedConcurrentChannel[nodeHasPeersForHash]
	scrape                       concurrency.BufferedConcurrentChannel[nodeHasPeersForHash]
	requestMetaInfo              concurrency.BufferedConcurrentChannel[infoHashWithPeers]
	persistTorrents              concurrency.BatchingChannel[infoHashWithMetaInfo]
	persistSources               concurrency.BatchingChannel[infoHashWithScrape]
	rescrapeThreshold            time.Duration
	saveFilesThreshold           uint
	savePieces                   bool
	dao                          *dao.Query
	classifierPublisher          publisher.Publisher[message.ClassifyTorrentPayload]
	// ignoreHashes is a thread-safe bloom filter that the crawler keeps in memory, containing every hash it has already encountered.
	// This avoids multiple attempts to crawl the same hash, and takes a lot of load off the database query that checks if a hash
	// has already been indexed. It is cleared every 6 hours.
	ignoreHashes *ignoreHashes
	// soughtNodeID is a random node ID used as the target for find_node and sample_infohashes requests.
	// It is rotated every 10 seconds.
	soughtNodeID *concurrency.AtomicValue[protocol.ID]
	stopped      chan struct{}
	logger       *zap.SugaredLogger
}

func (c *crawler) start() {
	// wait for the server to be ready
	select {
	case <-c.stopped:
		return
	case <-c.client.Ready():
		break
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// start the various pipeline workers
	go c.rotateSoughtNodeId(ctx)
	go c.runDiscoveredNodes(ctx)
	go c.runPing(ctx)
	go c.runFindNode(ctx)
	go c.getNodesForFindNode(ctx)
	go c.runSampleInfoHashes(ctx)
	go c.getNodesForSampleInfoHashes(ctx)
	go c.runInfoHashTriage(ctx)
	go c.runGetPeers(ctx)
	go c.runRequestMetaInfo(ctx)
	go c.runScrape(ctx)
	go c.reseedBootstrapNodes(ctx)
	go c.runPersistTorrents(ctx)
	go c.runPersistSources(ctx)
	go c.getOldNodes(ctx)
	go c.rotateIgnoreHashes(ctx)
	<-c.stopped
}

type nodeHasPeersForHash struct {
	infoHash protocol.ID
	node     netip.AddrPort
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
	bloom bloom.Filter
}

func (i *ignoreHashes) testOrAdd(id protocol.ID) bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.bloom.TestOrAdd(id[:])
}

func (i *ignoreHashes) clearAll() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.bloom.ClearAll()
}

func (c *crawler) rotateIgnoreHashes(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(6 * time.Hour):
			c.ignoreHashes.clearAll()
		}
	}
}

func (c *crawler) rotateSoughtNodeId(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(10 * time.Second):
			c.soughtNodeID.Set(protocol.RandomNodeID())
		}
	}
}
