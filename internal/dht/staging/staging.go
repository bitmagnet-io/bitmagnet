package staging

import (
	"context"
	"errors"
	"github.com/anacrolix/dht/v2/krpc"
	"github.com/bitmagnet-io/bitmagnet/internal/bloom"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/message"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/publisher"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
	"sync"
	"time"
)

type Params struct {
	fx.In
	Lifecycle           fx.Lifecycle
	Search              search.Search
	Dao                 *dao.Query
	ClassifierPublisher publisher.Publisher[message.ClassifyTorrentPayload]
	Logger              *zap.SugaredLogger
	Config              dht.Config
}

type Result struct {
	fx.Out
	Staging Staging
}

type Staging interface {
	Start() error
	Shutdown() error
	Stage(hashes ...InfoHashWithPeer)
	Requested() <-chan Request
	Reject(hash model.Hash20)
	Respond(ctx context.Context, response Response)
	Count() uint
}

func New(p Params) (r Result, err error) {
	r.Staging = &staging{
		mutex:               &sync.RWMutex{},
		activeRequests:      make(map[model.Hash20]Request),
		holding:             make(chan InfoHashWithPeer),
		holdingSize:         100,
		maxHoldingTime:      time.Second,
		saveFiles:           p.Config.SaveFiles,
		savePieces:          p.Config.SavePieces,
		rescrapeThreshold:   p.Config.RescrapeThreshold,
		requested:           make(chan Request),
		search:              p.Search,
		dao:                 p.Dao,
		classifierPublisher: p.ClassifierPublisher,
		logger:              p.Logger.Named("dht_staging"),
	}
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go (func() {
				// todo Handle error
				_ = r.Staging.Start()
			})()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return r.Staging.Shutdown()
		},
	})
	return
}

type InfoHashWithPeer struct {
	InfoHash model.Hash20
	Peer     krpc.NodeAddr
}

type Request struct {
	InfoHashWithPeer
	NeedMetaInfo       bool
	NeedScrape         bool
	torrentIsPersisted bool
}

type Response struct {
	InfoHash model.Hash20
	MetaInfo metainfo.Info
	Scrape   ResponseScrape
}

type ResponseScrape struct {
	Bfsd    bloom.Filter
	Bfpe    bloom.Filter
	Scraped bool
}

type staging struct {
	mutex               *sync.RWMutex
	started             bool
	stop                func()
	activeRequests      map[model.Hash20]Request
	holding             chan InfoHashWithPeer
	holdingSize         uint
	maxHoldingTime      time.Duration
	rescrapeThreshold   time.Duration
	saveFiles           bool
	savePieces          bool
	requested           chan Request
	search              search.TorrentSearch
	dao                 *dao.Query
	classifierPublisher publisher.Publisher[message.ClassifyTorrentPayload]
	logger              *zap.SugaredLogger
}

var (
	ErrorAlreadyStarted = errors.New("staging already started")
	ErrorStagingStopped = errors.New("staging stopped")
)

func (s *staging) Start() error {
	s.mutex.Lock()
	unlockOnce := sync.Once{}
	unlock := func() {
		unlockOnce.Do(s.mutex.Unlock)
	}
	defer unlock()
	if s.started {
		return ErrorAlreadyStarted
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.stop = cancel
	s.started = true
	unlock()
	go s.awaitHoldingHashes(ctx)
	<-ctx.Done()
	if stopErr := s.Shutdown(); stopErr != nil {
		return stopErr
	}
	return ErrorStagingStopped
}

func (s *staging) Shutdown() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.started {
		s.started = false
		s.stop()
	}
	return nil
}

func (s *staging) Stage(hashes ...InfoHashWithPeer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, h := range hashes {
		if _, ok := s.activeRequests[h.InfoHash]; !ok {
			s.holding <- h
		}
	}
}

func (s *staging) awaitHoldingHashes(ctx context.Context) {
	holdingMutex := &sync.Mutex{}
	holdingHashes := make(map[model.Hash20]InfoHashWithPeer, s.holdingSize)
	flushLocked := func() {
		if len(holdingHashes) == 0 {
			return
		}
		requests := s.getIndexerRequests(ctx, holdingHashes)
		s.mutex.Lock()
		for _, req := range requests {
			s.activeRequests[req.InfoHash] = req
		}
		s.mutex.Unlock()
		for _, req := range requests {
			s.requested <- req
		}
		holdingHashes = make(map[model.Hash20]InfoHashWithPeer, s.holdingSize)
	}
	flush := func() {
		holdingMutex.Lock()
		defer holdingMutex.Unlock()
		flushLocked()
	}
	add := func(h InfoHashWithPeer) {
		holdingMutex.Lock()
		defer holdingMutex.Unlock()
		s.mutex.RLock()
		_, isActive := s.activeRequests[h.InfoHash]
		s.mutex.RUnlock()
		if isActive {
			return
		}
		holdingHashes[h.InfoHash] = h
		if len(holdingHashes) >= int(s.holdingSize) {
			flushLocked()
		}
	}
	for {
		select {
		case <-ctx.Done():
			return
		case h, ok := <-s.holding:
			if !ok {
				return
			}
			go add(h)
		case <-time.After(s.maxHoldingTime):
			go flush()
		}
	}
}

func (s *staging) getIndexerRequests(ctx context.Context, hashesWithPeers map[model.Hash20]InfoHashWithPeer) []Request {
	hashes := make([]model.Hash20, 0, len(hashesWithPeers))
	for _, h := range hashesWithPeers {
		hashes = append(hashes, h.InfoHash)
	}
	searchResult, searchErr := s.search.Torrents(
		ctx,
		query.Where(search.TorrentInfoHashCriteria(hashes...)),
		query.Preload(func(q *dao.Query) []field.RelationField {
			return []field.RelationField{
				q.Torrent.Sources.RelationField,
			}
		}),
	)
	if searchErr != nil {
		s.logger.Errorf("failed to search existing torrents: %s", searchErr.Error())
		return nil
	}
	foundTorrents := make(map[model.Hash20]model.Torrent)
	for _, t := range searchResult.Items {
		foundTorrents[t.InfoHash] = t
	}
	requests := make([]Request, 0, len(hashesWithPeers))
	for _, h := range hashesWithPeers {
		if t, ok := foundTorrents[h.InfoHash]; !ok {
			requests = append(requests, Request{
				InfoHashWithPeer:   h,
				NeedMetaInfo:       true,
				NeedScrape:         true,
				torrentIsPersisted: false,
			})
		} else {
			needMetaInfo := s.saveFiles && !t.HasFilesInfo()
			needScrape := true
			for _, src := range t.Sources {
				if src.Source == "dht" {
					if src.Seeders.Valid && src.Leechers.Valid && src.UpdatedAt.After(time.Now().Add(-s.rescrapeThreshold)) {
						needScrape = false
					}
					break
				}
			}
			if needMetaInfo || needScrape {
				requests = append(requests, Request{
					InfoHashWithPeer:   h,
					NeedMetaInfo:       needMetaInfo,
					NeedScrape:         needScrape,
					torrentIsPersisted: true,
				})
			}
		}
	}
	return requests
}

func (s *staging) Requested() <-chan Request {
	return s.requested
}

func (s *staging) Reject(hash model.Hash20) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.activeRequests, hash)
}

func (s *staging) Respond(ctx context.Context, res Response) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	req, ok := s.activeRequests[res.InfoHash]
	if !ok {
		s.logger.Errorf("response did not correspond to an active request: %#v", res.InfoHash.String())
		return
	}
	defer delete(s.activeRequests, res.InfoHash)
	if req.NeedMetaInfo && res.MetaInfo.PieceLength > 0 {
		t, tErr := createTorrentModel(res.InfoHash, res.MetaInfo, res.Scrape)
		if tErr != nil {
			s.logger.Errorf("error creating torrent model: %s", tErr.Error())
			return
		}
		if !s.saveFiles && t.SingleFile.Valid && !t.SingleFile.Bool {
			t.SingleFile = model.NullBool{}
			t.Files = nil
		}
		if !s.savePieces {
			t.PieceLength = model.NullUint64{}
			t.Pieces = nil
		}
		putTorrentErr := s.dao.WithContext(ctx).Torrent.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&t)
		if putTorrentErr != nil {
			s.logger.Errorf("error persisting torrent: %s", putTorrentErr.Error())
			return
		}
		_, classifyErr := s.classifierPublisher.Publish(ctx, message.ClassifyTorrentPayload{
			InfoHashes: []model.Hash20{t.InfoHash},
		})
		if classifyErr != nil {
			s.logger.Errorf("error publishing classify message: %s", classifyErr.Error())
			return
		}
	} else if req.torrentIsPersisted && req.NeedScrape && res.Scrape.Scraped {
		src, srcErr := createTorrentSourceModel(res.InfoHash, res.Scrape)
		if srcErr != nil {
			s.logger.Errorf("error creating torrent source model: %s", srcErr.Error())
			return
		}
		putSourceErr := s.dao.WithContext(ctx).TorrentsTorrentSource.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&src)
		if putSourceErr != nil {
			s.logger.Errorf("error persisting torrent source: %s", putSourceErr.Error())
			return
		}
	}
}

func (s *staging) Count() uint {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return uint(len(s.activeRequests))
}
