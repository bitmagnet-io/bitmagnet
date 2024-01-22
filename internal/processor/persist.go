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
	return c.dao.Transaction(func(tx *dao.Query) error {
		torrentContentsPtr := make([]*model.TorrentContent, 0, len(torrentContents))
		deleteHashes := make([]driver.Valuer, 0, len(torrentContents))
		for _, content := range torrentContents {
			c := content
			if c.ContentID.Valid {
				deleteHashes = append(deleteHashes, c.InfoHash)
			}
			c.Torrent = model.Torrent{}
			torrentContentsPtr = append(torrentContentsPtr, &c)
		}
		if len(deleteHashes) > 0 {
			if _, deleteErr := tx.TorrentContent.WithContext(ctx).Where(
				c.dao.TorrentContent.InfoHash.In(deleteHashes...),
				//).Scopes(func(d gen.Dao) gen.Dao {
				//	return d.Where(c.dao.TorrentContent.ContentType.IsNull()).Or(c.dao.TorrentContent.ContentID.IsNull())
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
