package search

import (
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const LanguageFacetKey = "language"

func TorrentContentLanguageFacet(options ...query.FacetOption) query.Facet {
	return torrentContentLanguageFacet{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(LanguageFacetKey),
				query.FacetHasLabel("Language"),
				query.FacetUsesLogic(model.FacetLogicOr),
			}, options...)...,
		),
	}
}

type torrentContentLanguageFacet struct {
	query.FacetConfig
}

func (f torrentContentLanguageFacet) Aggregate(ctx query.FacetContext) (query.AggregationItems, error) {
	var results []struct {
		Language model.Language
		Count    uint
	}
	q, qErr := ctx.NewAggregationQuery()
	if qErr != nil {
		return nil, qErr
	}
	tx := q.UnderlyingDB().Select(
		"jsonb_array_elements(torrent_contents.languages) as language",
		"count(*) as count",
	).Group(
		"language",
	).Find(&results)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to aggregate languages: %w", tx.Error)
	}
	agg := make(query.AggregationItems, len(results))
	for _, item := range results {
		agg[item.Language.Id()] = query.AggregationItem{
			Label: item.Language.Name(),
			Count: item.Count,
		}
	}
	return agg, nil
}

func (f torrentContentLanguageFacet) Criteria() []query.Criteria {
	return []query.Criteria{
		query.GenCriteria(func(ctx query.DbContext) (query.Criteria, error) {
			filter := f.Filter().Values()
			langs := make([]model.Language, 0, len(filter))
			for _, v := range filter {
				lang := model.ParseLanguage(v)
				if !lang.Valid {
					return nil, errors.New("invalid language filter specified")
				}
				langs = append(langs, lang.Language)
			}
			if len(langs) == 0 {
				return query.AndCriteria{}, nil
			}
			array := "array["
			for i, lang := range langs {
				if i > 0 {
					array += ","
				}
				array += fmt.Sprintf("'%s'", lang.Id())
			}
			array += "]"
			return query.RawCriteria{
				Query: "torrent_contents.languages ?| " + array,
				Joins: maps.NewInsertMap(maps.MapEntry[string, struct{}]{Key: model.TableNameTorrentContent}),
			}, nil
		}),
	}
}
