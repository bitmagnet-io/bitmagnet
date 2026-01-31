package persister

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
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
	job := &persistJob{
		flusher: f,
		payload: payload.flatten(),
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
	startTime time.Time
	stats     AllTablesStats
}

func (j *persistJob) run(ctx context.Context) error {
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
		j.startTime = time.Now()

		if err := j.handleTorrentSources(ctx, tx); err != nil {
			return err
		}

		if err := j.handleTorrents(ctx, tx); err != nil {
			return err
		}

		missingInfoHashes, err := j.missingInfoHashes(ctx, tx)
		if err != nil {
			return err
		}

		if err := j.handleTorrentPieces(ctx, tx, missingInfoHashes); err != nil {
			return err
		}

		if err := j.handleTorrentsTorrentSources(ctx, tx, missingInfoHashes); err != nil {
			return err
		}

		if err := j.handleTorrentFiles(ctx, tx, missingInfoHashes); err != nil {
			return err
		}

		if err := j.handleContent(ctx, tx); err != nil {
			return err
		}

		if err := j.handleTorrentContent(ctx, tx, missingInfoHashes); err != nil {
			return err
		}

		if err := j.handleTorrentTags(ctx, tx, missingInfoHashes); err != nil {
			return err
		}

		if err := j.handleDeleteInfoHashes(ctx, tx); err != nil {
			return err
		}

		if err := j.handleQueueJobs(ctx, tx); err != nil {
			return err
		}

		return nil
	})
}

func (j *persistJob) handleTorrentSources(ctx context.Context, tx *dao.Query) error {
	if j.torrentSources.Len() == 0 {
		return nil
	}

	var stats TableStats

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

	stats.Affected += int(result.RowsAffected)

	for _, m := range torrentSourcesPtr {
		if m.CreatedAt.Before(j.startTime) {
			stats.Updated++
		} else {
			stats.Created++
		}
	}

	j.stats.Add(tx.TorrentSource.TableName(), stats)

	return nil
}

func (j *persistJob) handleTorrents(ctx context.Context, tx *dao.Query) error {
	if j.torrents.Len() == 0 {
		return nil
	}

	var stats TableStats

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

	stats.Affected += int(result.RowsAffected)

	for _, m := range torrentsPtr {
		if m.CreatedAt.Before(j.startTime) {
			stats.Updated++
		} else {
			stats.Created++
		}
	}

	j.stats.Add(tx.Torrent.TableName(), stats)

	return nil
}

func (j *persistJob) handleTorrentPieces(
	ctx context.Context,
	tx *dao.Query,
	missingInfoHashes map[protocol.ID]struct{},
) error {
	var stats TableStats

	if torrentPieces := slice.Filter(
		j.torrentPieces.Values(),
		func(m model.TorrentPieces) bool {
			if _, ok := missingInfoHashes[m.InfoHash]; ok || j.deleteInfoHashes.Has(m.InfoHash) {
				stats.Ignored++
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

		stats.Affected += int(result.RowsAffected)

		for _, m := range torrentPiecesPtr {
			if m.CreatedAt.Before(j.startTime) {
				stats.Ignored++
			} else {
				stats.Created++
			}
		}
	}

	j.stats.Add(tx.TorrentPieces.TableName(), stats)

	return nil
}

func (j *persistJob) handleTorrentsTorrentSources(
	ctx context.Context,
	tx *dao.Query,
	missingInfoHashes map[protocol.ID]struct{},
) error {
	var stats TableStats

	if torrentsTorrentSources := slice.Filter(
		j.torrentsTorrentSources.Values(),
		func(m model.TorrentsTorrentSource) bool {
			if _, ok := missingInfoHashes[m.InfoHash]; ok || j.deleteInfoHashes.Has(m.InfoHash) {
				stats.Ignored++
				return false
			}

			return true
		},
	); len(torrentsTorrentSources) > 0 {
		torrentsTorrentSourcesPtr := sliceToPointers(torrentsTorrentSources)

		result := tx.TorrentsTorrentSource.WithContext(ctx).
			Clauses(
				clause.Returning{
					Columns: []clause.Column{
						{Name: string(tx.TorrentsTorrentSource.CreatedAt.ColumnName())},
					},
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

		stats.Affected += int(result.RowsAffected)

		for _, t := range torrentsTorrentSourcesPtr {
			if t.CreatedAt.Before(j.startTime) {
				stats.Updated++
			} else {
				stats.Created++
			}
		}
	}

	j.stats.Add(tx.TorrentsTorrentSource.TableName(), stats)

	return nil
}

func (j *persistJob) handleTorrentFiles(
	ctx context.Context,
	tx *dao.Query,
	missingInfoHashes map[protocol.ID]struct{},
) error {
	var stats TableStats

	if torrentFiles := slice.Filter(
		j.torrentFiles.Values(),
		func(m model.TorrentFile) bool {
			if _, ok := missingInfoHashes[m.InfoHash]; ok || j.deleteInfoHashes.Has(m.InfoHash) {
				stats.Ignored++
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

		stats.Affected += int(result.RowsAffected)

		for _, m := range torrentFilesPtr {
			if m.CreatedAt.Before(j.startTime) {
				stats.Updated++
			} else {
				stats.Created++
			}
		}
	}

	j.stats.Add(tx.TorrentFile.TableName(), stats)

	return nil
}

func (j *persistJob) handleContent(ctx context.Context, tx *dao.Query) error {
	if j.content.Len() == 0 {
		return nil
	}

	var stats TableStats

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

	stats.Affected += int(result.RowsAffected)

	for _, m := range contentPtr {
		if m.CreatedAt.Before(j.startTime) {
			stats.Updated++
		} else {
			stats.Created++
		}
	}

	j.stats.Add(tx.Content.TableName(), stats)

	return nil
}

func (j *persistJob) handleTorrentContent(
	ctx context.Context,
	tx *dao.Query,
	missingInfoHashes map[protocol.ID]struct{},
) error {
	var stats TableStats

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

		stats.Affected += int(result.RowsAffected)
		stats.Deleted += int(result.RowsAffected)
	}

	if torrentContents := slice.Filter(
		j.torrentContents.Values(),
		func(m model.TorrentContent) bool {
			if _, ok := missingInfoHashes[m.InfoHash]; ok || j.deleteInfoHashes.Has(m.InfoHash) {
				stats.Ignored++
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

		stats.Affected += int(result.RowsAffected)

		for _, m := range torrentContentsPtr {
			if m.CreatedAt.Before(j.startTime) {
				stats.Updated++
			} else {
				stats.Created++
			}
		}
	}

	j.stats.Add(tx.TorrentContent.TableName(), stats)

	return nil
}

func (j *persistJob) handleTorrentTags(
	ctx context.Context,
	tx *dao.Query,
	missingInfoHashes map[protocol.ID]struct{},
) error {
	var (
		stats       TableStats
		torrentTags []*model.TorrentTag
	)

	torrentTagsToDelete := maps.NewInsertMap[hashWithID, struct{}]()

	for _, entry := range j.torrentTags.Entries() {
		if _, ok := missingInfoHashes[entry.Key.hash]; ok {
			if entry.Value {
				stats.Ignored++
			}

			continue
		}

		if entry.Value {
			torrentTags = append(torrentTags, &model.TorrentTag{
				InfoHash: entry.Key.hash,
				Name:     entry.Key.id,
			})
		} else {
			torrentTagsToDelete.SetKey(entry.Key)
		}
	}

	if len(torrentTags) > 0 {
		result := tx.TorrentTag.WithContext(ctx).Clauses(
			clause.Returning{
				Columns: []clause.Column{{Name: string(tx.TorrentTag.CreatedAt.ColumnName())}},
			},
			clause.OnConflict{
				DoNothing: true,
			},
		).
			UnderlyingDB().
			CreateInBatches(torrentTags, 100)
		if result.Error != nil {
			return result.Error
		}

		stats.Affected += int(result.RowsAffected)

		for _, m := range torrentTags {
			if m.CreatedAt.Before(j.startTime) {
				stats.Updated++
			} else {
				stats.Created++
			}
		}
	}

	// todo: Batch delete?
	for _, key := range torrentTagsToDelete.Keys() {
		result, err := tx.TorrentTag.WithContext(ctx).Where(
			tx.TorrentTag.InfoHash.Eq(key.hash),
			tx.TorrentTag.Name.Eq(key.id),
		).Delete()
		if err != nil {
			return err
		}

		stats.Deleted += int(result.RowsAffected)
	}

	j.stats.Add(tx.TorrentTag.TableName(), stats)

	return nil
}

func (j *persistJob) handleDeleteInfoHashes(ctx context.Context, tx *dao.Query) error {
	if j.deleteInfoHashes.Len() == 0 {
		return nil
	}

	var stats TableStats

	valuers := slice.Map(j.deleteInfoHashes.Keys(), func(infoHash protocol.ID) driver.Valuer {
		return infoHash
	})

	result, err := tx.Torrent.WithContext(ctx).Where(
		tx.Torrent.InfoHash.In(valuers...),
	).Delete()
	if err != nil {
		return err
	}

	stats.Affected += int(result.RowsAffected)
	stats.Deleted += int(result.RowsAffected)

	j.stats.Add(tx.Torrent.TableName(), stats)

	return nil
}

func (j *persistJob) handleQueueJobs(ctx context.Context, tx *dao.Query) error {
	if j.queueJobs.Len() == 0 {
		return nil
	}

	var stats TableStats

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

	stats.Affected += int(result.RowsAffected)
	stats.Created += int(result.RowsAffected)

	j.stats.Add(tx.QueueJob.TableName(), stats)

	return nil
}

func sliceToPointers[T any](sl []T) []*T {
	return slice.Map(sl, asPointer)
}

func asPointer[T any](v T) *T {
	return &v
}
