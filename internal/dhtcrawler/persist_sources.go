package dhtcrawler

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gen"
	"gorm.io/gorm/clause"
)

// runPersistSources waits on the persistSources channel for scraped torrents, and persists sources
// (which includes discovery date, seeders and leechers) to the database in batches.
func (cr *crawler) runPersistSources(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
	shutdown := make(chan struct{})

	go func() {
		defer cancel(nil)

		for {
			select {
			case <-ctx.Done():
				return

			case <-shutdown:
				return

			case scrapes := <-cr.persistSources.Out():
				srcs := make([]*model.TorrentsTorrentSource, 0, len(scrapes))

				hashSet := make(map[protocol.ID]struct{}, len(scrapes))
				for _, s := range scrapes {
					if _, ok := hashSet[s.infoHash]; ok {
						continue
					}

					hashSet[s.infoHash] = struct{}{}

					if src, err := createTorrentSourceModel(s); err != nil {
						cr.logger.Errorf("error creating torrent source model: %s", err.Error())
					} else {
						srcs = append(srcs, &src)
					}
				}

				dao, err := cr.daoProvider.Dao()
				if err != nil {
					cr.logger.Errorf("failed to acquire database: %s", err.Error())
					break
				}

				if persistErr := dao.WithContext(ctx).TorrentsTorrentSource.Clauses(
					clause.OnConflict{
						Columns: []clause.Column{
							{Name: string(dao.TorrentsTorrentSource.InfoHash.ColumnName())},
							{Name: string(dao.TorrentsTorrentSource.Source.ColumnName())},
						},
						DoUpdates: clause.AssignmentColumns([]string{
							string(dao.TorrentsTorrentSource.Seeders.ColumnName()),
							string(dao.TorrentsTorrentSource.Leechers.ColumnName()),
							// sets to null, fixes torrents indexed before 0.8.0 with
							// published_at
							// 0001-01-01 00:00:00+00:
							string(dao.TorrentsTorrentSource.PublishedAt.ColumnName()),
							string(dao.TorrentsTorrentSource.UpdatedAt.ColumnName()),
						}),
					},
				).Where(
					// check that the torrent record hasn't been deleted:
					gen.Exists(dao.WithContext(ctx).Torrent.Where(
						dao.Torrent.InfoHash.EqCol(dao.TorrentsTorrentSource.InfoHash),
					)),
				).CreateInBatches(srcs, 100); persistErr != nil {
					cr.logger.Errorf("error persisting torrent sources: %s", persistErr.Error())
				} else {
					cr.persistedTotal.With(prometheus.Labels{"entity": "TorrentsTorrentSource"}).Add(float64(len(srcs)))
					cr.logger.Debugw("persisted torrent sources", "count", len(srcs))
				}
			}
		}
	}()

	return func(context.Context) error {
		close(shutdown)

		<-ctx.Done()

		return nil
	}, nil
}

func createTorrentSourceModel(
	result infoHashWithScrape,
) (model.TorrentsTorrentSource, error) {
	seeders := model.NewNullUint(uint(result.bfsd.ApproximatedSize()))
	leechers := model.NewNullUint(uint(result.bfpe.ApproximatedSize()))

	return model.TorrentsTorrentSource{
		Source:   "dht",
		InfoHash: result.infoHash,
		Seeders:  seeders,
		Leechers: leechers,
	}, nil
}
