package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

func TorrentContentTypeCriteria(types ...model.ContentType) query.Criteria {
	return query.DaoCriteria{
		Conditions: func(ctx query.DbContext) ([]field.Expr, error) {
			q := ctx.Query()
			conditions := make([]field.Expr, 0, len(types))
			for _, contentType := range types {
				conditions = append(conditions, q.TorrentContent.ContentType.Eq(contentType))
			}
			return conditions, nil
		},
		Joins: maps.NewInsertMap(
			maps.MapEntry[string, struct{}]{Key: model.TableNameTorrentContent},
		),
	}
}
