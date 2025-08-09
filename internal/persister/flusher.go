package persister

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"gorm.io/gorm/clause"
)

type flusher struct {
	daoProvider database.DaoTransactionProvider
	blocker     blocker.Blocker
	sem         chan struct{}
}

func (f *flusher) flush(ctx context.Context, payload *payload) (AllTablesStats, error) {
	job := persistJob{
		flusher: f,
		payload: *payload,
		stats:   make(AllTablesStats),
	}

	err := job.run(ctx)
	if err != nil {
		return nil, err
	}

	return job.stats, nil
}

type persistJob struct {
	*flusher
	payload
	stats AllTablesStats
}

func (j persistJob) run(ctx context.Context) error {
	if j.len() == 0 {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case j.sem <- struct{}{}:
	}

	defer func() { <-j.sem }()

	if j.deleteInfoHashes.Len() > 0 {
		if blockErr := j.blocker.Block(ctx, j.deleteInfoHashes.Keys(), false); blockErr != nil {
			return blockErr
		}
	}

	return j.daoProvider.DaoTransaction(func(tx *dao.Query) error {
		startTime := time.Now()

		var torrentSourceStats TableStats

		if j.torrentSources.Len() > 0 {
			torrentSourcesPtr := sliceToPointers(j.torrentSources.Values())

			result := tx.TorrentSource.WithContext(ctx).Clauses(
				clause.Returning{
					Columns: []clause.Column{{Name: string(tx.TorrentSource.CreatedAt.ColumnName())}},
				},
				clause.OnConflict{
					UpdateAll: true,
				}).
				UnderlyingDB().
				CreateInBatches(torrentSourcesPtr, 100)
			if result.Error != nil {
				return result.Error
			}

			torrentSourceStats.Affected += int(result.RowsAffected)
			for _, m := range torrentSourcesPtr {
				if m.CreatedAt.Before(startTime) {
					torrentSourceStats.Updated++
				} else {
					torrentSourceStats.Created++
				}
			}
		}

		var torrentStats TableStats

		if j.torrents.Len() > 0 {
			// todo: check for all scenerios!
			torrentsPtr := sliceToPointers(j.torrents.Values())

			result := tx.Torrent.WithContext(ctx).Clauses(
				clause.Returning{
					Columns: []clause.Column{{Name: string(tx.Torrent.CreatedAt.ColumnName())}},
				},
				clause.OnConflict{
					Columns: []clause.Column{{Name: string(tx.Torrent.InfoHash.ColumnName())}},
					DoUpdates: clause.AssignmentColumns([]string{
						string(tx.Torrent.Name.ColumnName()),
						string(tx.Torrent.FilesStatus.ColumnName()),
						string(tx.Torrent.FilesCount.ColumnName()),
						string(tx.Torrent.UpdatedAt.ColumnName()),
					}),
				}).
				UnderlyingDB().
				CreateInBatches(torrentsPtr, 100)

			if result.Error != nil {
				return result.Error
			}

			torrentStats.Affected += int(result.RowsAffected)

			for _, m := range torrentsPtr {
				if m.CreatedAt.Before(startTime) {
					torrentStats.Updated++
				} else {
					torrentStats.Created++
				}
			}
		}

		missingInfoHashes, err := j.payload.missingInfoHashes(ctx, tx)
		if err != nil {
			return err
		}

		var torrentPiecesStats TableStats

		if torrentPieces := slice.Filter(
			j.torrentPieces.Values(),
			func(m model.TorrentPieces) bool {
				if _, ok := missingInfoHashes[m.InfoHash]; ok || j.payload.deleteInfoHashes.Has(m.InfoHash) {
					torrentPiecesStats.Ignored++
					return false
				}

				return true
			},
		); len(torrentPieces) > 0 {
			torrentPiecesPtr := sliceToPointers(torrentPieces)

			result := tx.TorrentPieces.WithContext(ctx).Clauses(
				clause.Returning{
					Columns: []clause.Column{{Name: string(tx.TorrentPieces.CreatedAt.ColumnName())}},
				},
				clause.OnConflict{
					DoNothing: true,
				}).
				UnderlyingDB().
				CreateInBatches(torrentPiecesPtr, 100)
			if result.Error != nil {
				return result.Error
			}

			torrentPiecesStats.Affected += int(result.RowsAffected)

			for _, m := range torrentPiecesPtr {
				if m.CreatedAt.Before(startTime) {
					torrentPiecesStats.Ignored++
				} else {
					torrentPiecesStats.Created++
				}
			}
		}

		j.stats.Add(tx.TorrentPieces.TableName(), torrentPiecesStats)

		var torrentsTorrentSourcesStats TableStats

		if torrentsTorrentSources := slice.Filter(
			j.torrentsTorrentSources.Values(),
			func(m model.TorrentsTorrentSource) bool {
				if _, ok := missingInfoHashes[m.InfoHash]; ok || j.payload.deleteInfoHashes.Has(m.InfoHash) {
					torrentsTorrentSourcesStats.Ignored++
					return false
				}

				return true
			},
		); len(torrentsTorrentSources) > 0 {
			torrentsTorrentSourcesPtr := sliceToPointers(torrentsTorrentSources)

			result := tx.TorrentsTorrentSource.WithContext(ctx).
				Clauses(
					clause.Returning{
						Columns: []clause.Column{{Name: string(tx.TorrentsTorrentSource.CreatedAt.ColumnName())}},
					},
					clause.OnConflict{
						Columns: []clause.Column{
							{Name: string(tx.TorrentsTorrentSource.InfoHash.ColumnName())},
							{Name: string(tx.TorrentsTorrentSource.Source.ColumnName())},
						},
						DoUpdates: clause.AssignmentColumns([]string{
							string(tx.TorrentsTorrentSource.Seeders.ColumnName()),
							string(tx.TorrentsTorrentSource.Leechers.ColumnName()),
							string(tx.TorrentsTorrentSource.UpdatedAt.ColumnName()),
						}),
					},
				).
				UnderlyingDB().
				CreateInBatches(torrentsTorrentSourcesPtr, 100)
			if result.Error != nil {
				return result.Error
			}

			torrentsTorrentSourcesStats.Affected += int(result.RowsAffected)

			for _, t := range torrentsTorrentSourcesPtr {
				if t.CreatedAt.Before(startTime) {
					torrentsTorrentSourcesStats.Updated++
				} else {
					torrentsTorrentSourcesStats.Created++
				}
			}
		}

		j.stats.Add(tx.TorrentsTorrentSource.TableName(), torrentsTorrentSourcesStats)

		var torrentFilesStats TableStats

		if torrentFiles := slice.Filter(
			j.torrentFiles.Values(),
			func(m model.TorrentFile) bool {
				if _, ok := missingInfoHashes[m.InfoHash]; ok || j.payload.deleteInfoHashes.Has(m.InfoHash) {
					torrentFilesStats.Ignored++
					return false
				}

				return true
			},
		); len(torrentFiles) > 0 {
			torrentFilesPtr := sliceToPointers(torrentFiles)

			result := tx.TorrentFile.WithContext(ctx).Clauses(
				clause.Returning{
					Columns: []clause.Column{{Name: string(tx.TorrentFile.CreatedAt.ColumnName())}},
				},
				clause.OnConflict{
					UpdateAll: true,
				}).
				UnderlyingDB().
				CreateInBatches(torrentFilesPtr, 100)
			if result.Error != nil {
				return result.Error
			}

			torrentFilesStats.Affected += int(result.RowsAffected)

			for _, m := range torrentFilesPtr {
				if m.CreatedAt.Before(startTime) {
					torrentFilesStats.Updated++
				} else {
					torrentFilesStats.Created++
				}
			}
		}

		j.stats.Add(tx.TorrentFile.TableName(), torrentFilesStats)

		var contentStats TableStats

		if j.content.Len() > 0 {
			contentPtr := sliceToPointers(j.content.Values())

			result := tx.Content.WithContext(ctx).Clauses(
				clause.Returning{
					Columns: []clause.Column{{Name: string(tx.Content.CreatedAt.ColumnName())}},
				},
				clause.OnConflict{
					UpdateAll: true,
				}).
				UnderlyingDB().
				CreateInBatches(contentPtr, 100)
			if result.Error != nil {
				return result.Error
			}

			contentStats.Affected += int(result.RowsAffected)

			for _, m := range contentPtr {
				if m.CreatedAt.Before(startTime) {
					contentStats.Updated++
				} else {
					contentStats.Created++
				}
			}
		}

		j.stats.Add(tx.Content.TableName(), contentStats)

		var torrentContentStats TableStats

		if j.deleteTorrentContent.Len() > 0 {
			result, err := tx.TorrentContent.WithContext(ctx).Where(
				tx.TorrentContent.ID.In(slice.Map(
					j.deleteTorrentContent.Keys(),
					func(ref model.TorrentContentRef) string {
						return ref.InferID()
					},
				)...),
			).Delete()
			if err != nil {
				return err
			}

			torrentContentStats.Affected += int(result.RowsAffected)
			torrentContentStats.Deleted += int(result.RowsAffected)
		}

		if torrentContents := slice.Filter(
			j.torrentContents.Values(),
			func(m model.TorrentContent) bool {
				if _, ok := missingInfoHashes[m.InfoHash]; ok || j.payload.deleteInfoHashes.Has(m.InfoHash) {
					torrentContentStats.Ignored++
					return false
				}

				return true
			},
		); len(torrentContents) > 0 {
			torrentContentsPtr := sliceToPointers(torrentContents)

			result := tx.TorrentContent.WithContext(ctx).Clauses(
				clause.Returning{
					Columns: []clause.Column{{Name: string(tx.TorrentContent.CreatedAt.ColumnName())}},
				},
				clause.OnConflict{
					UpdateAll: true,
				},
			).
				UnderlyingDB().
				CreateInBatches(torrentContentsPtr, 100)
			if result.Error != nil {
				return result.Error
			}

			torrentContentStats.Affected += int(result.RowsAffected)
			for _, m := range torrentContentsPtr {
				if m.CreatedAt.Before(startTime) {
					torrentContentStats.Updated++
				} else {
					torrentContentStats.Created++
				}
			}
		}

		j.stats.Add(tx.TorrentContent.TableName(), torrentContentStats)

		var torrentTagsStats TableStats

		if torrentTags := slice.Filter(
			j.torrentTags.Values(),
			func(m model.TorrentTag) bool {
				if _, ok := missingInfoHashes[m.InfoHash]; ok || j.payload.deleteInfoHashes.Has(m.InfoHash) {
					torrentTagsStats.Ignored++
					return false
				}

				return true
			},
		); len(torrentTags) > 0 {
			torrentTagsPtr := sliceToPointers(torrentTags)

			result := tx.TorrentTag.WithContext(ctx).Clauses(
				clause.Returning{
					Columns: []clause.Column{{Name: string(tx.TorrentTag.CreatedAt.ColumnName())}},
				},
				clause.OnConflict{
					DoNothing: true,
				},
			).
				UnderlyingDB().
				CreateInBatches(torrentTagsPtr, 100)
			if result.Error != nil {
				return result.Error
			}

			torrentTagsStats.Affected += int(result.RowsAffected)

			for _, m := range torrentTagsPtr {
				if m.CreatedAt.Before(startTime) {
					torrentTagsStats.Updated++
				} else {
					torrentTagsStats.Created++
				}
			}
		}

		j.stats.Add(tx.TorrentTag.TableName(), torrentTagsStats)

		if j.deleteInfoHashes.Len() > 0 {
			valuers := slice.Map(j.deleteInfoHashes.Keys(), func(infoHash protocol.ID) driver.Valuer {
				return infoHash
			})

			result, err := tx.Torrent.WithContext(ctx).Where(
				tx.Torrent.InfoHash.In(valuers...),
			).Delete()
			if err != nil {
				return err
			}

			torrentStats.Affected += int(result.RowsAffected)
			torrentStats.Deleted += int(result.RowsAffected)
		}

		j.stats.Add(tx.Torrent.TableName(), torrentStats)

		var queueJobStats TableStats

		if j.queueJobs.Len() > 0 {
			result := tx.QueueJob.WithContext(ctx).Clauses(
				clause.OnConflict{
					DoNothing: true,
				},
			).UnderlyingDB().
				CreateInBatches(slice.Map(j.queueJobs.Values(), func(j model.QueueJob) *model.QueueJob {
					return &j
				}), 100)
			if result.Error != nil {
				return result.Error
			}

			queueJobStats.Affected += int(result.RowsAffected)
			queueJobStats.Created += int(result.RowsAffected)
		}

		return nil
	})
}

func sliceToPointers[T any](sl []T) []*T {
	return slice.Map(sl, asPointer)
}

func asPointer[T any](v T) *T {
	return &v
}
