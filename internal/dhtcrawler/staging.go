package dhtcrawler

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/bloom"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/message"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/publisher"
	"go.uber.org/zap"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
	"net/netip"
	"sync"
	"time"
)

type infoHashWithPeer struct {
	infoHash protocol.ID
	peer     netip.AddrPort
}

type stagingRequest struct {
	infoHashWithPeer
	needMetaInfo       bool
	needScrape         bool
	torrentIsPersisted bool
}

type stagingResponse struct {
	infoHash protocol.ID
	metaInfo metainfo.Info
	scrape   stagingResponseScrape
}

type stagingResponseScrape struct {
	bfsd    bloom.Filter
	bfpe    bloom.Filter
	scraped bool
}

type staging struct {
	mutex                  sync.RWMutex
	activeRequests         map[protocol.ID]stagingRequest
	requestHolding         chan infoHashWithPeer
	requestHoldingSize     uint
	maxRequestHoldingTime  time.Duration
	responseHolding        chan stagingResponse
	responseHoldingSize    uint
	maxResponseHoldingTime time.Duration
	rescrapeThreshold      time.Duration
	saveFiles              bool
	saveFilesThreshold     uint
	savePieces             bool
	requested              chan stagingRequest
	search                 search.TorrentSearch
	dao                    *dao.Query
	classifierPublisher    publisher.Publisher[message.ClassifyTorrentPayload]
	stopped                chan struct{}
	logger                 *zap.SugaredLogger
}

func (s *staging) start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go s.awaitHoldingHashes(ctx)
	<-s.stopped
}

func (s *staging) stage(hashes ...infoHashWithPeer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, h := range hashes {
		if _, ok := s.activeRequests[h.infoHash]; !ok {
			s.requestHolding <- h
		}
	}
}

func (s *staging) awaitHoldingHashes(ctx context.Context) {
	ticker := time.NewTicker(s.maxRequestHoldingTime)
	holdingMutex := &sync.Mutex{}
	holdingHashes := make(map[protocol.ID]infoHashWithPeer, s.requestHoldingSize)
	flushLocked := func() {
		if len(holdingHashes) == 0 {
			return
		}
		requests := s.getIndexerRequests(ctx, holdingHashes)
		s.mutex.Lock()
		for _, req := range requests {
			s.activeRequests[req.infoHash] = req
		}
		s.mutex.Unlock()
		for _, req := range requests {
			s.requested <- req
		}
		holdingHashes = make(map[protocol.ID]infoHashWithPeer, s.requestHoldingSize)
	}
	flush := func() {
		holdingMutex.Lock()
		defer holdingMutex.Unlock()
		flushLocked()
	}
	add := func(h infoHashWithPeer) {
		holdingMutex.Lock()
		defer holdingMutex.Unlock()
		s.mutex.RLock()
		_, isActive := s.activeRequests[h.infoHash]
		s.mutex.RUnlock()
		if isActive {
			return
		}
		if len(holdingHashes) == 0 {
			ticker.Reset(s.maxRequestHoldingTime)
		}
		holdingHashes[h.infoHash] = h
		if len(holdingHashes) >= int(s.requestHoldingSize) {
			flushLocked()
		}
	}
	for {
		select {
		case <-ctx.Done():
			return
		case h := <-s.requestHolding:
			go add(h)
		case <-ticker.C:
			go flush()
		}
	}
}

func (s *staging) getIndexerRequests(ctx context.Context, hashesWithPeers map[protocol.ID]infoHashWithPeer) []stagingRequest {
	hashes := make([]protocol.ID, 0, len(hashesWithPeers))
	for _, h := range hashesWithPeers {
		hashes = append(hashes, h.infoHash)
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
	foundTorrents := make(map[protocol.ID]model.Torrent)
	for _, t := range searchResult.Items {
		foundTorrents[t.InfoHash] = t
	}
	requests := make([]stagingRequest, 0, len(hashesWithPeers))
	for _, h := range hashesWithPeers {
		if t, ok := foundTorrents[h.infoHash]; !ok {
			requests = append(requests, stagingRequest{
				infoHashWithPeer:   h,
				needMetaInfo:       true,
				needScrape:         true,
				torrentIsPersisted: false,
			})
		} else {
			needMetaInfo := s.saveFiles && t.WantFilesInfo()
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
				requests = append(requests, stagingRequest{
					infoHashWithPeer:   h,
					needMetaInfo:       needMetaInfo,
					needScrape:         needScrape,
					torrentIsPersisted: true,
				})
			}
		}
	}
	return requests
}

func (s *staging) drop(hash protocol.ID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.activeRequests, hash)
}

func (s *staging) awaitResponses(ctx context.Context) {
	ticker := time.NewTicker(s.maxResponseHoldingTime)
	holdingMutex := &sync.Mutex{}
	holdingResponses := make(map[protocol.ID]stagingResponse, s.responseHoldingSize)
	flushLocked := func() {
		if len(holdingResponses) == 0 {
			return
		}
		responses := make([]stagingResponse, 0, len(holdingResponses))
		for _, res := range holdingResponses {
			responses = append(responses, res)
		}
		s.handleResponseBatch(ctx, responses...)
		holdingResponses = make(map[protocol.ID]stagingResponse, s.responseHoldingSize)
	}
	flush := func() {
		holdingMutex.Lock()
		defer holdingMutex.Unlock()
		flushLocked()
	}
	add := func(r stagingResponse) {
		holdingMutex.Lock()
		defer holdingMutex.Unlock()
		if len(holdingResponses) == 0 {
			ticker.Reset(s.maxResponseHoldingTime)
		}
		holdingResponses[r.infoHash] = r
		if len(holdingResponses) >= int(s.responseHoldingSize) {
			flushLocked()
		}
	}
	for {
		select {
		case <-ctx.Done():
			return
		case r := <-s.responseHolding:
			go add(r)
		case <-ticker.C:
			go flush()
		}
	}
}

func (s *staging) handleResponseBatch(ctx context.Context, responses ...stagingResponse) {
	hashesHandled := make(map[protocol.ID]struct{})
	torrentsToPersist := make([]*model.Torrent, 0, len(responses))
	sourcesToPersist := make([]*model.TorrentsTorrentSource, 0)
	hashesToClassify := make([]protocol.ID, 0, len(responses))
	s.mutex.Lock()
	for _, res := range responses {
		req, ok := s.activeRequests[res.infoHash]
		if !ok {
			s.logger.Errorf("response did not correspond to an active request: %#v", res.infoHash)
			continue
		}
		hashesHandled[res.infoHash] = struct{}{}
		if req.needMetaInfo && res.metaInfo.PieceLength > 0 {
			t, tErr := createTorrentModel(res.infoHash, res.metaInfo, res.scrape)
			if tErr != nil {
				s.logger.Errorf("error creating torrent model: %s", tErr.Error())
				continue
			}
			if !s.saveFiles {
				if t.FilesStatus != model.FilesStatusSingle {
					t.FilesStatus = model.FilesStatusNoInfo
					t.Files = nil
				}
			} else if s.saveFilesThreshold > 0 && len(t.Files) > int(s.saveFilesThreshold) {
				t.FilesStatus = model.FilesStatusOverThreshold
				t.Files = nil
			}
			if !s.savePieces {
				t.PieceLength = model.NullUint64{}
				t.Pieces = nil
			}
			torrentsToPersist = append(torrentsToPersist, &t)
			hashesToClassify = append(hashesToClassify, t.InfoHash)
		} else if req.torrentIsPersisted && req.needScrape && res.scrape.scraped {
			src, srcErr := createTorrentSourceModel(res.infoHash, res.scrape)
			if srcErr != nil {
				s.logger.Errorf("error creating torrent source model: %s", srcErr.Error())
				continue
			}
			sourcesToPersist = append(sourcesToPersist, &src)
		}
	}
	// give the mutex a breather while we persist everything...
	s.mutex.Unlock()
	if len(torrentsToPersist) > 0 {
		if persistErr := s.dao.WithContext(ctx).Torrent.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).CreateInBatches(torrentsToPersist, 100); persistErr != nil {
			s.logger.Errorf("error persisting torrents: %s", persistErr.Error())
		} else {
			s.logger.Debugf("persisted %d torrents", len(torrentsToPersist))
			if _, classifyErr := s.classifierPublisher.Publish(ctx, message.ClassifyTorrentPayload{
				InfoHashes: hashesToClassify,
			}); classifyErr != nil {
				s.logger.Errorf("error publishing classify message: %s", classifyErr.Error())
			}
		}
	}
	if len(sourcesToPersist) > 0 {
		if persistErr := s.dao.WithContext(ctx).TorrentsTorrentSource.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).CreateInBatches(sourcesToPersist, 100); persistErr != nil {
			s.logger.Errorf("error persisting torrent sources: %s", persistErr.Error())
		}
	}
	s.mutex.Lock()
	for h := range hashesHandled {
		delete(s.activeRequests, h)
	}
	s.mutex.Unlock()
}

func (s *staging) count() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.activeRequests)
}
