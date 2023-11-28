package dhtcrawler

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/message"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm/clause"
)

// runPersistTorrents waits on the persistTorrents channel, and persists torrents to the database in batches.
// After persisting each batch it will publish a message to the classifier,
// and forward the hash on the scrape channel to attempt finding the seeders/leechers.
func (c *crawler) runPersistTorrents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case is := <-c.persistTorrents.Out():
			ts := make([]*model.Torrent, 0, len(is))
			hashMap := make(map[protocol.ID]infoHashWithMetaInfo, len(is))
			for _, i := range is {
				if _, ok := hashMap[i.infoHash]; ok {
					continue
				}
				hashMap[i.infoHash] = i
				if t, err := createTorrentModel(i.infoHash, i.metaInfo, c.savePieces, c.saveFilesThreshold); err != nil {
					c.logger.Errorf("error creating torrent model: %s", err.Error())
				} else {
					ts = append(ts, &t)
				}
			}
			if persistErr := c.dao.WithContext(ctx).Torrent.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: string(c.dao.Torrent.InfoHash.ColumnName())}},
				DoUpdates: clause.AssignmentColumns([]string{
					string(c.dao.Torrent.Name.ColumnName()),
					string(c.dao.Torrent.FilesStatus.ColumnName()),
					string(c.dao.Torrent.PieceLength.ColumnName()),
					string(c.dao.Torrent.Pieces.ColumnName()),
				}),
			}).CreateInBatches(ts, 100); persistErr != nil {
				c.logger.Errorf("error persisting torrents: %s", persistErr)
			} else {
				c.persistedTotal.With(prometheus.Labels{"entity": "Torrent"}).Add(float64(len(ts)))
				c.logger.Debugw("persisted torrents", "count", len(ts))
				hashesToClassify := make([]protocol.ID, 0, classifyBatchSize)
				flushClassify := func() {
					if len(hashesToClassify) == 0 {
						return
					}
					if _, classifyErr := c.classifierPublisher.Publish(ctx, message.ClassifyTorrentPayload{
						InfoHashes: hashesToClassify,
					}); classifyErr != nil {
						c.logger.Errorf("error publishing classify message: %s", classifyErr.Error())
					}
				}
				for _, t := range ts {
					hashesToClassify = append(hashesToClassify, t.InfoHash)
					if len(hashesToClassify) == classifyBatchSize {
						flushClassify()
						hashesToClassify = make([]protocol.ID, 0, classifyBatchSize)
					}
				}
				flushClassify()
				for _, i := range hashMap {
					select {
					case <-ctx.Done():
						return
					case c.scrape.In() <- i.nodeHasPeersForHash:
						continue
					}
				}
			}
		}
	}
}

func createTorrentModel(
	hash protocol.ID,
	info metainfo.Info,
	savePieces bool,
	saveFilesThreshold uint,
) (model.Torrent, error) {
	name := info.BestName()
	private := false
	if info.Private != nil {
		private = *info.Private
	}
	var files []model.TorrentFile
	for i, file := range info.Files {
		files = append(files, model.TorrentFile{
			Index: uint32(i),
			Path:  file.DisplayPath(&info),
			Size:  uint64(file.Length),
		})
	}
	filesStatus := model.FilesStatusSingle
	if len(files) > int(saveFilesThreshold) {
		filesStatus = model.FilesStatusOverThreshold
		files = nil
	} else if len(files) > 0 {
		filesStatus = model.FilesStatusMulti
	}
	var pieceLength model.NullUint64
	var pieces []byte
	if savePieces {
		pieceLength = model.NewNullUint64(uint64(info.PieceLength))
		pieces = info.Pieces
	}
	return model.Torrent{
		InfoHash:    hash,
		Name:        name,
		Size:        uint64(info.TotalLength()),
		Private:     private,
		PieceLength: pieceLength,
		Pieces:      pieces,
		Files:       files,
		FilesStatus: filesStatus,
		Sources: []model.TorrentsTorrentSource{
			{
				Source:   "dht",
				InfoHash: hash,
			},
		},
	}, nil
}

const classifyBatchSize = 200

// runPersistSources waits on the persistSources channel for scraped torrents, and persists sources
// (which includes discovery date, seeders and leechers) to the database in batches.
func (c *crawler) runPersistSources(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case scrapes := <-c.persistSources.Out():
			srcs := make([]*model.TorrentsTorrentSource, 0, len(scrapes))
			hashSet := make(map[protocol.ID]struct{}, len(scrapes))
			for _, s := range scrapes {
				if _, ok := hashSet[s.infoHash]; ok {
					continue
				}
				hashSet[s.infoHash] = struct{}{}
				if src, err := createTorrentSourceModel(s); err != nil {
					c.logger.Errorf("error creating torrent source model: %s", err.Error())
				} else {
					srcs = append(srcs, &src)
				}
			}
			if persistErr := c.dao.WithContext(ctx).TorrentsTorrentSource.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).CreateInBatches(srcs, 100); persistErr != nil {
				c.logger.Errorf("error persisting torrent sources: %s", persistErr.Error())
			} else {
				c.persistedTotal.With(prometheus.Labels{"entity": "TorrentsTorrentSource"}).Add(float64(len(srcs)))
				c.logger.Debugw("persisted torrent sources", "count", len(srcs))
			}
		}
	}
}

func createTorrentSourceModel(
	result infoHashWithScrape,
) (model.TorrentsTorrentSource, error) {
	seeders := model.NewNullUint(uint(result.bfsd.ApproximatedSize()))
	leechers := model.NewNullUint(uint(result.bfpe.ApproximatedSize()))
	// todo add discovered result to bloom?
	bfsdBytes, bfsdErr := result.bfsd.MarshalBinary()
	if bfsdErr != nil {
		return model.TorrentsTorrentSource{}, bfsdErr
	}
	bfpeBytes, bfpeErr := result.bfpe.MarshalBinary()
	if bfpeErr != nil {
		return model.TorrentsTorrentSource{}, bfpeErr
	}
	return model.TorrentsTorrentSource{
		Source:   "dht",
		InfoHash: result.infoHash,
		Bfsd:     bfsdBytes,
		Bfpe:     bfpeBytes,
		Leechers: seeders,
		Seeders:  leechers,
	}, nil
}
