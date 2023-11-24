package search

import (
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

const TorrentTagFacetKey = "torrent_tag"

func TorrentTagsFacet(options ...query.FacetOption) query.Facet {
	return torrentTagFacet{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(TorrentTagFacetKey),
				query.FacetHasLabel("Torrent Tag"),
				query.FacetUsesAndLogic(),
				query.FacetHasAggregationOption(query.RequireJoin(model.TableNameTorrentContent)),
			}, options...)...,
		),
	}
}

type torrentTagFacet struct {
	query.FacetConfig
}

func (f torrentTagFacet) Aggregate(ctx query.FacetContext) (query.AggregationItems, error) {
	var results []struct {
		Value string
		Count uint
	}
	q, qErr := ctx.NewAggregationQuery(
		query.Table(model.TableNameTorrentTag),
		query.Join(func(daoQ *dao.Query) []query.TableJoin {
			return []query.TableJoin{
				{
					Table: daoQ.TorrentContent,
					On: []field.Expr{
						daoQ.TorrentContent.InfoHash.EqCol(daoQ.TorrentTag.InfoHash),
					},
					Type: query.TableJoinTypeInner,
				},
			}
		}),
	)
	if qErr != nil {
		return nil, qErr
	}
	if err := q.UnderlyingDB().Select(
		fmt.Sprintf("%s.name as value", model.TableNameTorrentTag),
		"count(*) as count",
	).Group(
		"value",
	).Find(&results).Error; err != nil {
		return nil, err
	}
	agg := make(query.AggregationItems, len(results))
	var values []string
	for _, item := range results {
		agg[item.Value] = query.AggregationItem{
			Count: item.Count,
		}
		values = append(values, item.Value)
	}
	if len(values) > 0 {
		tags, tagsErr := ctx.Query().TorrentTag.WithContext(ctx.Context()).Where(
			ctx.Query().TorrentTag.Name.In(values...),
		).Find()
		if tagsErr != nil {
			return nil, tagsErr
		}
		for _, tag := range tags {
			thisAgg := agg[tag.Name]
			thisAgg.Label = tag.Name
			agg[tag.Name] = thisAgg
		}
	}
	return agg, nil
}

func (f torrentTagFacet) Criteria() []query.Criteria {
	filter := f.Filter().Values()
	criteria := make([]query.Criteria, len(filter))
	for i, tag := range filter {
		criteria[i] = TorrentTagCriteria(tag)
	}
	return criteria
}
