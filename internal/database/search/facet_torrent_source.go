package search

import (
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

const TorrentSourceFacetKey = "torrent_source"

func TorrentSourceFacet(options ...query.FacetOption) query.Facet {
	return torrentSourceFacet{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(TorrentSourceFacetKey),
				query.FacetHasLabel("Torrent Source"),
				query.FacetUsesOrLogic(),
			}, options...)...,
		),
	}
}

type torrentSourceFacet struct {
	query.FacetConfig
}

func (f torrentSourceFacet) Aggregate(ctx query.FacetContext) (query.AggregationItems, error) {
	var results []struct {
		Value string
		Count uint
	}
	q, qErr := ctx.NewAggregationQuery(
		query.Table(model.TableNameTorrentsTorrentSource),
		query.Join(func(daoQ *dao.Query) []query.TableJoin {
			return []query.TableJoin{
				{
					Table: daoQ.Torrent,
					On: []field.Expr{
						daoQ.Torrent.InfoHash.EqCol(daoQ.TorrentsTorrentSource.InfoHash),
					},
					Type: query.TableJoinTypeInner,
				},
				{
					Table: daoQ.TorrentContent,
					On: []field.Expr{
						daoQ.TorrentContent.InfoHash.EqCol(daoQ.TorrentsTorrentSource.InfoHash),
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
		fmt.Sprintf("%s.source as value", model.TableNameTorrentsTorrentSource),
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
		sources, sourcesErr := ctx.Query().TorrentSource.WithContext(ctx.Context()).Where(
			ctx.Query().TorrentSource.Key.In(values...),
		).Find()
		if sourcesErr != nil {
			return nil, sourcesErr
		}
		for _, source := range sources {
			thisAgg := agg[source.Key]
			thisAgg.Label = source.Name
			agg[source.Key] = thisAgg
		}
	}
	return agg, nil
}

func (f torrentSourceFacet) Criteria() []query.Criteria {
	filter := f.Filter().Values()
	if len(filter) == 0 {
		return []query.Criteria{}
	}
	return []query.Criteria{
		TorrentSourceCriteria(filter...),
	}
}

func TorrentSourceCriteria(keys ...string) query.Criteria {
	return query.GenCriteria(func(ctx query.DbContext) (query.Criteria, error) {
		q := ctx.Query()
		return query.RawCriteria{
			Query: gen.Exists(
				q.TorrentsTorrentSource.Where(
					q.TorrentsTorrentSource.InfoHash.EqCol(q.TorrentContent.InfoHash),
					q.TorrentsTorrentSource.Source.In(keys...),
				),
			),
		}, nil
	})
}
