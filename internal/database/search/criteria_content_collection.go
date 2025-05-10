package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen"
)

type collectionMap = map[string]map[string]map[string]struct{}

func collectionMapFromRefs(refs ...model.ContentCollectionRef) collectionMap {
	m := make(collectionMap)
	for _, ref := range refs {
		if _, ok := m[ref.Type]; !ok {
			m[ref.Type] = make(map[string]map[string]struct{})
		}

		if _, ok := m[ref.Type][ref.Source]; !ok {
			m[ref.Type][ref.Source] = make(map[string]struct{})
		}

		m[ref.Type][ref.Source][ref.ID] = struct{}{}
	}

	return m
}

func ContentCollectionCriteria(refs ...model.ContentCollectionRef) query.Criteria {
	return query.GenCriteria(func(ctx query.DBContext) (query.Criteria, error) {
		refMap := collectionMapFromRefs(refs...)
		q := ctx.Query()

		var criteria []query.Criteria

		for collectionType, sourceMap := range refMap {
			for source, idMap := range sourceMap {
				ids := make([]string, 0, len(idMap))
				for id := range idMap {
					ids = append(ids, id)
				}

				criteria = append(criteria, query.RawCriteria{
					Joins: maps.NewInsertMap(
						maps.MapEntry[string, struct{}]{Key: q.TorrentContent.TableName()},
					),
					Query: gen.Exists(
						q.ContentCollectionContent.Where(
							q.ContentCollectionContent.ContentType.EqCol(
								q.TorrentContent.ContentType,
							),
							q.ContentCollectionContent.ContentSource.EqCol(
								q.TorrentContent.ContentSource,
							),
							q.ContentCollectionContent.ContentID.EqCol(
								q.TorrentContent.ContentID,
							),
							q.ContentCollectionContent.ContentCollectionType.Eq(
								collectionType,
							),
							q.ContentCollectionContent.ContentCollectionSource.Eq(source),
							q.ContentCollectionContent.ContentCollectionID.In(ids...),
						),
					),
				})
			}
		}

		return query.And(
			query.RawCriteria{
				Query: q.TorrentContent.ContentID.IsNotNull(),
			},
			query.Or(criteria...),
		), nil
	})
}
