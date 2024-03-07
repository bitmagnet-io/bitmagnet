package processor

import (
	"context"
	"database/sql/driver"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"gorm.io/gorm/clause"
)

type persistPayload struct {
	torrentContents  []model.TorrentContent
	deleteIds        []string
	deleteInfoHashes []protocol.ID
}

func (c processor) persist(ctx context.Context, payload persistPayload) error {
	contentsMap := make(map[model.ContentRef]struct{}, len(payload.torrentContents))
	contentsPtr := make([]*model.Content, 0, len(payload.torrentContents))
	torrentContentsPtr := make([]*model.TorrentContent, 0, len(payload.torrentContents))
	for _, tc := range payload.torrentContents {
		tcCopy := tc
		tcCopy.Torrent = model.Torrent{}
		if tcCopy.ContentID.Valid && tcCopy.Content.CreatedAt.IsZero() {
			contentRef := tcCopy.Content.Ref()
			if _, ok := contentsMap[contentRef]; !ok {
				contentsMap[contentRef] = struct{}{}
				contentCopy := tcCopy.Content
				contentsPtr = append(contentsPtr, &contentCopy)
			}
		}
		tcCopy.Content = model.Content{}
		torrentContentsPtr = append(torrentContentsPtr, &tcCopy)
	}
	// a semaphore is used here to avoid a Postgres deadlock being detected when multiple processes are trying to persist
	if err := c.persistSemaphore.Acquire(ctx, 1); err != nil {
		return err
	}
	defer c.persistSemaphore.Release(1)
	return c.dao.Transaction(func(tx *dao.Query) error {
		if len(contentsPtr) > 0 {
			if createContentErr := tx.Content.WithContext(ctx).Clauses(
				clause.OnConflict{
					UpdateAll: true,
				}).CreateInBatches(contentsPtr, 100); createContentErr != nil {
				return createContentErr
			}
		}
		if len(payload.deleteInfoHashes) > 0 {
			valuers := make([]driver.Valuer, 0, len(payload.deleteInfoHashes))
			for _, infoHash := range payload.deleteInfoHashes {
				valuers = append(valuers, infoHash)
			}
			if _, deleteErr := tx.Torrent.WithContext(ctx).Where(
				c.dao.Torrent.InfoHash.In(valuers...),
			).Delete(); deleteErr != nil {
				return deleteErr
			}
		}
		if len(payload.deleteIds) > 0 {
			if _, deleteErr := tx.TorrentContent.WithContext(ctx).Where(
				c.dao.TorrentContent.ID.In(payload.deleteIds...),
			).Delete(); deleteErr != nil {
				return deleteErr
			}
		}
		return tx.TorrentContent.WithContext(ctx).Clauses(
			clause.OnConflict{
				UpdateAll: true,
			},
		).CreateInBatches(torrentContentsPtr, 100)
	})
}
