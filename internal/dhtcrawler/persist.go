package dhtcrawler

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/asynq/message"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
	"gorm.io/gorm/clause"
)

func (c *crawler) runPersistTorrents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case is := <-c.persistTorrents.Out():
			ts := make([]*model.Torrent, 0, len(is))
			for _, i := range is {
				if t, err := createTorrentModel(i.infoHash, i.metaInfo, c.saveFilesThreshold); err != nil {
					c.logger.Errorf("error creating torrent model: %s", err.Error())
				} else {
					ts = append(ts, &t)
				}
			}
			if persistErr := c.dao.WithContext(ctx).Torrent.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).CreateInBatches(ts, 100); persistErr != nil {
				c.logger.Errorf("error persisting torrents: %s", persistErr)
			} else {
				c.logger.Warnf("persisted %d torrents", len(ts))
				hashesToClassify := make([]protocol.ID, 0, len(ts))
				for _, t := range ts {
					hashesToClassify = append(hashesToClassify, t.InfoHash)
				}
				if _, classifyErr := c.classifierPublisher.Publish(ctx, message.ClassifyTorrentPayload{
					InfoHashes: hashesToClassify,
				}); classifyErr != nil {
					c.logger.Errorf("error publishing classify message: %s", classifyErr.Error())
				}
				for _, i := range is {
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
	return model.Torrent{
		InfoHash:    hash,
		Name:        name,
		Size:        uint64(info.TotalLength()),
		Private:     private,
		PieceLength: model.NewNullUint64(uint64(info.PieceLength)),
		Pieces:      info.Pieces,
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

func (c *crawler) runPersistSources(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case scrapes := <-c.persistSources.Out():
			srcs := make([]*model.TorrentsTorrentSource, 0, len(scrapes))
			for _, s := range scrapes {
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
				c.logger.Warnf("persisted %d torrent sources", len(srcs))
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
