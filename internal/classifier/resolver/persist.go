package resolver

import (
	"context"
	"database/sql/driver"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gorm/clause"
)

func (r resolver) Persist(ctx context.Context, torrentContents ...model.TorrentContent) error {
	if len(torrentContents) == 0 {
		return nil
	}
	return r.dao.Transaction(func(tx *dao.Query) error {
		torrentContentsPtr := make([]*model.TorrentContent, 0, len(torrentContents))
		deleteHashes := make([]driver.Valuer, 0, len(torrentContents))
		for _, content := range torrentContents {
			c := content
			if c.ContentID.Valid {
				deleteHashes = append(deleteHashes, c.InfoHash)
			}
			torrentContentsPtr = append(torrentContentsPtr, &c)
		}
		if len(deleteHashes) > 0 {
			if _, deleteErr := tx.TorrentContent.WithContext(ctx).Where(
				r.dao.TorrentContent.InfoHash.In(deleteHashes...),
				r.dao.TorrentContent.ContentID.IsNull(),
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
