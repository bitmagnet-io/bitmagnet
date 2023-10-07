package crawler

import (
	"context"
	"errors"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/worker"
	"github.com/bitmagnet-io/bitmagnet/internal/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/dht/routingtable"
	"github.com/bitmagnet-io/bitmagnet/internal/dht/server"
	"github.com/bitmagnet-io/bitmagnet/internal/dht/staging"
	"github.com/bitmagnet-io/bitmagnet/internal/metainfo/metainforequester"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Params struct {
	fx.In
	PeerID            krpc.ID `name:"dht_peer_id"`
	Config            dht.Config
	DhtServer         server.Server
	RoutingTable      routingtable.Table
	Staging           staging.Staging
	MetainfoRequester metainforequester.Requester
	Logger            *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Crawler Crawler
	Worker  worker.Worker `group:"workers"`
}

type Crawler interface {
	StartCrawling() error
	StopCrawling() error
}

func New(p Params) (r Result, err error) {
	r.Crawler = &crawler{
		peerID:                      p.PeerID,
		dhtServer:                   p.DhtServer,
		staging:                     p.Staging,
		routingTable:                p.RoutingTable,
		metainfoRequester:           p.MetainfoRequester,
		crawlBootstrapHostsInterval: p.Config.CrawlBootstrapHostsInterval,
		sampleInfoHashesInterval:    p.Config.SampleInfoHashesInterval,
		discardUnscrapableTorrents:  p.Config.DiscardUnscrapableTorrents,
		maxStagingSize:              p.Config.MaxStagingSize,
		savePieces:                  p.Config.SavePieces,
		mutex:                       &sync.Mutex{},
		logger:                      p.Logger.Named("dht_indexer"),
	}
	r.Worker = worker.NewWorker(
		"dht_crawler",
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				// todo Handle errors
				go (func() {
					_ = r.Crawler.StartCrawling()
				})()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return r.Crawler.StopCrawling()
			},
		},
	)
	return
}

type crawler struct {
	peerID                      [20]byte
	routingTable                routingtable.Table
	staging                     staging.Staging
	dhtServer                   server.Server
	metainfoRequester           metainforequester.Requester
	crawlBootstrapHostsInterval time.Duration
	sampleInfoHashesInterval    time.Duration
	discardUnscrapableTorrents  bool
	mutex                       *sync.Mutex
	crawling                    bool
	stopCrawling                func()
	maxStagingSize              uint
	savePieces                  bool
	logger                      *zap.SugaredLogger
}

var (
	ErrorAlreadyCrawling = errors.New("dht crawler: Already crawling")
	ErrorCrawlingStopped = errors.New("dht crawler: Crawling stopped")
)

func (c *crawler) StartCrawling() error {
	c.mutex.Lock()
	unlockOnce := sync.Once{}
	unlock := func() {
		unlockOnce.Do(c.mutex.Unlock)
	}
	defer unlock()
	if c.crawling {
		return ErrorAlreadyCrawling
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c.stopCrawling = cancel
	c.crawling = true
	unlock()
	go c.crawl(ctx)
	<-ctx.Done()
	if stopErr := c.StopCrawling(); stopErr != nil {
		return stopErr
	}
	return ErrorCrawlingStopped
}

func (c *crawler) crawl(ctx context.Context) {
	go c.crawlBootstrapHosts(ctx)
	go c.crawlPeersForInfoHashes(ctx)
	go c.awaitInfoHashes(ctx)
}

func (c *crawler) StopCrawling() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.crawling {
		c.crawling = false
		c.stopCrawling()
	}
	return nil
}
