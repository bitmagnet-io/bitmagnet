package processor

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gorm/clause"
)

func (c processor) persist(ctx context.Context, torrentContents []model.TorrentContent, deleteIds []string) error {
	contentsMap := make(map[model.ContentRef]struct{}, len(torrentContents))
	contentsPtr := make([]*model.Content, 0, len(torrentContents))
	torrentContentsPtr := make([]*model.TorrentContent, 0, len(torrentContents))
	for _, tc := range torrentContents {
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
		if len(deleteIds) > 0 {
			if _, deleteErr := tx.TorrentContent.WithContext(ctx).Where(
				c.dao.TorrentContent.ID.In(deleteIds...),
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
