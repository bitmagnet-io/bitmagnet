package dhtcrawler

import (
	"context"
	"database/sql/driver"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"time"
)

type triageResult struct {
	InfoHash    protocol.ID
	FilesStatus model.FilesStatus
	Seeders     model.NullInt
	Leechers    model.NullInt
	UpdatedAt   time.Time
}

func (c *crawler) runInfoHashTriage(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case reqs := <-c.infoHashTriage.Out():
			hashes := make([]driver.Valuer, 0, len(reqs))
			for _, r := range reqs {
				hashes = append(hashes, r.infoHash)
			}
			var result []*triageResult
			if queryErr := c.dao.Torrent.WithContext(ctx).Select(
				c.dao.Torrent.InfoHash,
				c.dao.Torrent.FilesStatus,
				c.dao.TorrentsTorrentSource.Seeders,
				c.dao.TorrentsTorrentSource.Leechers,
				c.dao.TorrentsTorrentSource.UpdatedAt,
			).LeftJoin(
				c.dao.TorrentsTorrentSource,
				c.dao.Torrent.InfoHash.EqCol(c.dao.TorrentsTorrentSource.InfoHash),
				c.dao.TorrentsTorrentSource.Source.Eq("dht"),
			).Where(
				c.dao.Torrent.InfoHash.In(hashes...),
			).UnderlyingDB().Find(&result).Error; queryErr != nil {
				c.logger.Errorf("failed to search existing torrents: %s", queryErr.Error())
				return
			}
			foundTorrents := make(map[protocol.ID]triageResult)
			for _, t := range result {
				foundTorrents[t.InfoHash] = *t
			}
			for _, r := range reqs {
				if t, ok := foundTorrents[r.infoHash]; !ok || t.FilesStatus == model.FilesStatusNoInfo {
					select {
					case <-ctx.Done():
						return
					case c.getPeers.In() <- r:
						continue
					}
				} else if !(t.Seeders.Valid && t.Leechers.Valid) || t.UpdatedAt.Before(time.Now().Add(-c.rescrapeThreshold)) {
					select {
					case <-ctx.Done():
						return
					case c.scrape.In() <- r:
						continue
					}
				}
			}
		}
	}
}
