package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen"
)

type contentMap = map[model.ContentType]map[string]map[string]struct{}

func contentMapFromRefs(refs ...model.ContentRef) contentMap {
	m := make(contentMap)
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

func ContentCanonicalIdentifierCriteria(refs ...model.ContentRef) query.Criteria {
	return contentCanonicalIdentifierCriteria(contentMapFromRefs(refs...))
}

func contentCanonicalIdentifierCriteria(contentMap contentMap) query.Criteria {
	return query.GenCriteria(func(ctx query.DbContext) (query.Criteria, error) {
		q := ctx.Query()
		criteria := make([]query.Criteria, 0, len(contentMap))
		for contentType, sourceMap := range contentMap {
			for source, idMap := range sourceMap {
				conds := make([]gen.Condition, 0, 3)
				if !contentType.IsNil() {
					conds = append(conds, q.Content.Type.Eq(contentType.String()))
				}
				ids := make([]string, 0, len(idMap))
				for id := range idMap {
					ids = append(ids, id)
				}
				conds = append(conds,
					q.Content.Source.Eq(source),
					q.Content.ID.In(ids...),
				)
				criteria = append(criteria, query.RawCriteria{
					Query: q.Content.Where(conds...).UnderlyingDB(),
				})
			}
		}
		return query.Or(criteria...), nil
	})
}

func ContentIdentifierCriteria(refs ...model.ContentRef) query.Criteria {
	return query.GenCriteria(func(ctx query.DbContext) (query.Criteria, error) {
		m := contentMapFromRefs(refs...)
		q := ctx.Query()
		criteria := []query.Criteria{contentCanonicalIdentifierCriteria(m)}
		for contentType, sourceMap := range m {
			for source, idMap := range sourceMap {
				conds := make([]gen.Condition, 0, 6)
				if !contentType.IsNil() {
					conds = append(conds, q.ContentAttribute.ContentType.Eq(contentType))
				}
				ids := make([]string, 0, len(idMap))
				for id := range idMap {
					ids = append(ids, id)
				}
				conds = append(conds,
					q.ContentAttribute.ContentType.EqCol(q.Content.Type),
					q.ContentAttribute.ContentSource.EqCol(q.Content.Source),
					q.ContentAttribute.ContentID.EqCol(q.Content.ID),
					q.ContentAttribute.Source.Eq(source),
					q.ContentAttribute.Value.In(ids...),
				)
				criteria = append(criteria, query.RawCriteria{
					Query: gen.Exists(
						q.ContentAttribute.Where(conds...),
					),
				})
			}
		}
		return query.OrCriteria{
			Criteria: criteria,
		}, nil
	})
}
