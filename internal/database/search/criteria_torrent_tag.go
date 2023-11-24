package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen"
)

func TorrentTagCriteria(tagNames ...string) query.Criteria {
	return query.GenCriteria(func(ctx query.DbContext) (query.Criteria, error) {
		q := ctx.Query()
		return query.OrCriteria{
			Criteria: []query.Criteria{
				query.RawCriteria{
					Query: gen.Exists(
						q.TorrentTag.Where(
							q.TorrentTag.InfoHash.EqCol(q.Torrent.InfoHash),
							q.TorrentTag.Name.In(tagNames...),
						),
					),
					Joins: maps.NewInsertMap(
						maps.MapEntry[string, struct{}]{Key: model.TableNameTorrent},
					),
				},
			},
		}, nil
	})
}
