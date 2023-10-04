package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
	"strings"
)

type contentCollectionFacet struct {
	query.FacetConfig
	collectionType string
}

func (r contentCollectionFacet) Aggregate(ctx query.FacetContext) (items query.AggregationItems, err error) {
	var results []struct {
		ConcatId string
		Name     string
		Count    uint
	}
	sq, sqErr := ctx.NewAggregationQuery(
		query.Table(model.TableNameContentCollectionContent),
		query.Join(func(q *dao.Query) []query.TableJoin {
			return []query.TableJoin{
				{
					Table: q.TorrentContent,
					On: []field.Expr{
						q.TorrentContent.ContentType.EqCol(q.ContentCollectionContent.ContentType),
						q.TorrentContent.ContentSource.EqCol(q.ContentCollectionContent.ContentSource),
						q.TorrentContent.ContentID.EqCol(q.ContentCollectionContent.ContentID),
					},
					Type: query.TableJoinTypeInner,
				},
			}
		}),
		query.RequireJoin(model.TableNameContentCollection),
		query.Where(query.RawCriteria{
			Query: "content_collections_content.content_collection_type = ?",
			Args:  []interface{}{r.collectionType},
		}),
	)
	if sqErr != nil {
		err = sqErr
		return
	}
	tx := sq.UnderlyingDB().Select(
		"(content_collections_content.content_collection_source || ':' ||content_collections_content.content_collection_id) as concat_id",
		"MIN(content_collections.name) as name",
		"count(distinct(content_collections_content.content_source, content_collections_content.content_id)) as count",
	).Group(
		"concat_id",
	).Find(&results)
	if tx.Error != nil {
		err = tx.Error
		return
	}
	agg := make(query.AggregationItems, len(results))
	for _, item := range results {
		agg[item.ConcatId] = query.AggregationItem{
			Label: item.Name,
			Count: item.Count,
		}
	}
	return agg, nil
}

func (r contentCollectionFacet) Criteria() []query.Criteria {
	sourceMap := make(map[string]map[string]struct{})
	filter := r.Filter().Values()
	for _, value := range filter {
		parts := strings.Split(value, ":")
		if len(parts) != 2 {
			continue
		}
		source, id := parts[0], parts[1]
		if _, ok := sourceMap[source]; !ok {
			sourceMap[source] = make(map[string]struct{})
		}
		sourceMap[source][id] = struct{}{}
	}
	criteria := make([]query.Criteria, 0, len(sourceMap))
	for source, idMap := range sourceMap {
		refs := make([]model.ContentCollectionRef, 0, len(idMap))
		for id := range idMap {
			refs = append(refs, model.ContentCollectionRef{
				Type:   r.collectionType,
				Source: source,
				ID:     id,
			})
		}
		switch r.Logic() {
		case model.FacetLogicOr:
			criteria = append(criteria, ContentCollectionCriteria(refs...))
		case model.FacetLogicAnd:
			for _, ref := range refs {
				criteria = append(criteria, ContentCollectionCriteria(ref))
			}
		}
	}
	return criteria
}
