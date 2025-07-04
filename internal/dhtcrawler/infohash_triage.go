package dhtcrawler

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

// runInfoHashTriage receives discovered hashes on the infoHashTriage channel, determines if they should be crawled,
// and forwards them to the appropriate channel. Possible outcomes are:
//  1. The hash is not in the database, so it is forwarded to the getPeers channel to attempt
//     retrieval of the meta info.
//  2. The hash is in the database, but we don't have the full details of the torrent (for example it was imported
//     outside the DHT crawler, and so we don't have the files info), so it is forwarded to the getPeers channel to
//     attempt retrieval of the meta info.
//  3. The hash is in the database, but the seeders/leechers are not known or are outdated,
//     so it is forwarded to the scrape channel.
//  4. The hash is in the database and the seeders/leechers are known and up to date, so it is discarded.
func (cr *crawler) runInfoHashTriage(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case reqs := <-cr.infoHashTriage.Out():
			allHashes := make([]protocol.ID, 0, len(reqs))

			reqMap := make(map[protocol.ID]nodeHasPeersForHash, len(reqs))
			for _, r := range reqs {
				if _, ok := reqMap[r.infoHash]; ok {
					continue
				}

				allHashes = append(allHashes, r.infoHash)
				reqMap[r.infoHash] = r
			}

			filteredHashes, filterErr := cr.blockingManager.Filter(ctx, allHashes)
			if filterErr != nil {
				cr.logger.Errorf("failed to filter infohashes: %s", filterErr.Error())
				break
			}

			if len(filteredHashes) == 0 {
				break
			}

			filteredHashMap := make(map[protocol.ID]struct{}, len(filteredHashes))
			valuers := make([]driver.Valuer, 0, len(filteredHashes))

			for _, h := range filteredHashes {
				filteredHashMap[h] = struct{}{}

				valuers = append(valuers, h)
			}

			dao, err := cr.daoProvider.Dao()
			if err != nil {
				cr.logger.Errorf("failed to acquire database: %s", err.Error())
				break
			}

			var result []*triageResult
			if queryErr := dao.Torrent.WithContext(ctx).Select(
				dao.Torrent.InfoHash,
				dao.Torrent.FilesStatus,
				dao.Torrent.FilesCount,
				dao.TorrentsTorrentSource.Seeders,
				dao.TorrentsTorrentSource.Leechers,
				dao.TorrentsTorrentSource.UpdatedAt,
			).LeftJoin(
				dao.TorrentsTorrentSource,
				dao.Torrent.InfoHash.EqCol(dao.TorrentsTorrentSource.InfoHash),
				dao.TorrentsTorrentSource.Source.Eq("dht"),
			).Where(
				dao.Torrent.InfoHash.In(valuers...),
			).UnderlyingDB().Find(&result).Error; queryErr != nil {
				cr.logger.Errorf("failed to search existing torrents: %s", queryErr.Error())
				break
			}

			foundTorrents := make(map[protocol.ID]triageResult)
			for _, t := range result {
				foundTorrents[t.InfoHash] = *t
			}

			for h := range filteredHashMap {
				r := reqMap[h]
				if t, ok := foundTorrents[r.infoHash]; !ok ||
					t.FilesStatus == model.FilesStatusNoInfo ||
					(t.FilesStatus != model.FilesStatusSingle && !t.FilesCount.Valid) ||
					(t.FilesStatus == model.FilesStatusOverThreshold && t.FilesCount.Uint <= cr.saveFilesThreshold) {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case cr.getPeers.In() <- r:
						continue
					}
				} else if (!t.Seeders.Valid || !t.Leechers.Valid) ||
					t.UpdatedAt.Before(time.Now().Add(-cr.rescrapeThreshold)) {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case cr.scrape.In() <- r:
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
	FilesCount  model.NullUint
	Seeders     model.NullUint
	Leechers    model.NullUint
	UpdatedAt   time.Time
}
