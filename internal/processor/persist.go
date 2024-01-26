package processor

import (
	"context"
	"database/sql/driver"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gorm/clause"
)

func (c processor) Persist(ctx context.Context, torrentContents ...model.TorrentContent) error {
	if len(torrentContents) == 0 {
		return nil
	}
	contentsMap := make(map[model.ContentRef]struct{}, len(torrentContents))
	contentsPtr := make([]*model.Content, 0, len(torrentContents))
	torrentContentsPtr := make([]*model.TorrentContent, 0, len(torrentContents))
	deleteHashes := make([]driver.Valuer, 0, len(torrentContents))
	for _, tc := range torrentContents {
		tcCopy := tc
		deleteHashes = append(deleteHashes, tcCopy.InfoHash)
		tcCopy.Torrent = model.Torrent{}
		if tcCopy.ContentID.Valid {
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
				}).CreateInBatches(contentsPtr, 20); createContentErr != nil {
				return createContentErr
			}
		}
		if _, deleteErr := tx.TorrentContent.WithContext(ctx).Where(
			c.dao.TorrentContent.InfoHash.In(deleteHashes...),
		).Delete(); deleteErr != nil {
			return deleteErr
		}
		return tx.TorrentContent.WithContext(ctx).Clauses(
			clause.OnConflict{
				DoNothing: true,
			},
		).CreateInBatches(torrentContentsPtr, 20)
	})
}
