package dhtcrawler

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/worker"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/message"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainforequester"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/publisher"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"net/netip"
	"sync"
	"time"
)

type Params struct {
	fx.In
	Config              Config
	KTable              ktable.TableBatch
	Server              server.Server
	MetainfoRequester   metainforequester.Requester
	Search              search.Search
	Dao                 *dao.Query
	ClassifierPublisher publisher.Publisher[message.ClassifyTorrentPayload]
	BootstrapNodes      []netip.AddrPort `name:"dht_bootstrap_nodes"`
	Logger              *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Worker worker.Worker `group:"workers"`
}

func New(params Params) Result {
	c := crawler{
		kTable: params.KTable,
		server: params.Server,
		staging: staging{
			mutex:                  sync.RWMutex{},
			activeRequests:         make(map[protocol.ID]stagingRequest),
			requestHolding:         make(chan infoHashWithPeer, 100),
			requestHoldingSize:     100,
			responseHolding:        make(chan stagingResponse, 2000),
			responseHoldingSize:    500,
			maxResponseHoldingTime: 10 * time.Second,
			maxRequestHoldingTime:  time.Second,
			saveFiles:              params.Config.SaveFiles,
			saveFilesThreshold:     params.Config.SaveFilesThreshold,
			savePieces:             params.Config.SavePieces,
			rescrapeThreshold:      params.Config.RescrapeThreshold,
			requested:              make(chan stagingRequest),
			search:                 params.Search,
			dao:                    params.Dao,
			classifierPublisher:    params.ClassifierPublisher,
			logger:                 params.Logger.Named("dht_staging"),
		},
		targetStagingSize:            1000,
		metainfoRequester:            params.MetainfoRequester,
		bootstrapNodes:               params.BootstrapNodes,
		reseedBootstrapNodesInterval: time.Minute * 10,
		getOldestPeersInterval:       time.Second * 10,
		oldPeerThreshold:             time.Minute * 15,
		findNodesInterval:            time.Second / 4,
		findNodeSemaphore:            semaphore.NewWeighted(1000),
		sampleInfoHashesSemaphore:    semaphore.NewWeighted(1000),
		sampleInfoHashesInterval:     time.Second,
		discoveredPeers: concurrency.NewBatchingDedupedChannel[ktable.Peer, string](
			make(chan ktable.Peer, 1000),
			100,
			time.Second/4,
			func(p ktable.Peer) string {
				return p.Addr().Addr().String()
			},
		),
		peersForPing:             newPeerChan(1000),
		peersForFindNode:         newPeerChan(1000),
		peersForSampleInfoHashes: newPeerChan(1000),
		soughtPeerId:             &concurrency.AtomicValue[protocol.ID]{},
		stopped:                  make(chan struct{}),
		logger:                   params.Logger.Named("dht_crawler"),
	}
	c.soughtPeerId.Set(protocol.RandomNodeID())
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
	kTable                       ktable.TableBatch
	server                       server.Server
	staging                      staging
	targetStagingSize            int
	metainfoRequester            metainforequester.Requester
	bootstrapNodes               []netip.AddrPort
	reseedBootstrapNodesInterval time.Duration
	getOldestPeersInterval       time.Duration
	oldPeerThreshold             time.Duration
	findNodesInterval            time.Duration
	findNodeSemaphore            *semaphore.Weighted
	sampleInfoHashesSemaphore    *semaphore.Weighted
	sampleInfoHashesInterval     time.Duration
	sampleInfoHashesShortfall    atomicValue[int]
	discoveredPeers              concurrency.BatchingDedupedChannel[ktable.Peer, string]
	peersForPing                 concurrency.BufferedDedupedChannel[ktable.Peer, string]
	peersForFindNode             concurrency.BufferedDedupedChannel[ktable.Peer, string]
	peersForSampleInfoHashes     concurrency.BufferedDedupedChannel[ktable.Peer, string]
	soughtPeerId                 *concurrency.AtomicValue[protocol.ID]
	stopped                      chan struct{}
	logger                       *zap.SugaredLogger
}

func newPeerChan(cap int) concurrency.BufferedDedupedChannel[ktable.Peer, string] {
	return concurrency.NewBufferedDedupedChannel[ktable.Peer, string](cap, func(p ktable.Peer) string {
		return p.Addr().Addr().String()
	})
}

func (c *crawler) start() {
	// wait for the server to be ready
	select {
	case <-c.stopped:
		return
	case <-c.server.Ready():
		break
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// start the various pipeline workers
	go c.rotateSoughtPeerId(ctx)
	go c.awaitDiscoveredPeers(ctx)
	go c.awaitPeersForPing(ctx)
	go c.awaitPeersForFindNode(ctx)
	go c.findNode(ctx)
	go c.awaitPeersForSampleInfoHashes(ctx)
	go c.sampleInfoHashes(ctx)
	// start seeding the table with the bootstrap nodes
	go c.reseedBootstrapNodes(ctx)
	// start the staging workers
	go c.staging.awaitHoldingHashes(ctx)
	go c.staging.awaitResponses(ctx)
	// await info hashes from staging
	go c.awaitInfoHashes(ctx)
	go c.getOldPeers(ctx)
	<-c.stopped
}

type atomicValue[T any] struct {
	mutex sync.RWMutex
	value T
}

func (a *atomicValue[T]) Get() T {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.value
}

func (a *atomicValue[T]) Set(value T) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.value = value
}

func (a *atomicValue[T]) Update(fn func(T) T) T {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.value = fn(a.value)
	return a.value
}
