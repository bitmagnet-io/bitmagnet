package dhtcrawler

import (
	"context"
	"database/sql/driver"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"time"
)

// runInfoHashTriage receives discovered hashes on the infoHashTriage channel, determines if they should be crawled,
// and forwards them to the appropriate channel. Possible outcomes are:
// 1. The hash is not in the database, so it is forwarded to the getPeers channel to attempt retrieval of the meta info.
// 2. The hash is in the database, but we don't have the full details of the torrent (for example it was imported outside the DHT crawler,
// and so we don't have the files info), so it is forwarded to the getPeers channel to attempt retrieval of the meta info.
// 3. The hash is in the database, but the seeders/leechers are not known or are outdated, so it is forwarded to the scrape channel.
// 4. The hash is in the database and the seeders/leechers are known and up to date, so it is discarded.
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
				break
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

type triageResult struct {
	InfoHash    protocol.ID
	FilesStatus model.FilesStatus
	Seeders     model.NullInt
	Leechers    model.NullInt
	UpdatedAt   time.Time
}
