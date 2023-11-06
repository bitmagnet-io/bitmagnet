package search

import (
	"database/sql/driver"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"gorm.io/gen/field"
)

func TorrentInfoHashCriteria(infoHashes ...protocol.ID) query.Criteria {
	valuers := make([]driver.Valuer, 0, len(infoHashes))
	for _, infoHash := range infoHashes {
		valuers = append(valuers, infoHash)
	}
	return query.DaoCriteria{
		Conditions: func(ctx query.DbContext) ([]field.Expr, error) {
			return []field.Expr{ctx.Query().Torrent.InfoHash.In(valuers...)}, nil
		},
	}
}
