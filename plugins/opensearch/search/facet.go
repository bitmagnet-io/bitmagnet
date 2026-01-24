//go:build wasip1

package search

import (
	"encoding/json"
	"errors"

	"github.com/bitmagnet-io/bitmagnet/proto/common/search"
	osquery "github.com/defensestation/osquery/v2"
)

type facet interface {
	Key() string
	queryAgg() osquery.Aggregation
	queryFilter(facet *search.FacetParam) osquery.Mappable
	resultAggItems(agg json.RawMessage) ([]*search.FacetResultItem, error)
}

type facets []facet

func (fs facets) req(req *osquery.SearchRequest, params *search.Params) error {
	aggsMap := make(map[string]struct{})
	filtersMap := make(map[string]osquery.Mappable)
	for _, param := range params.Facets {
		if f, ok := allFacetsMap[param.Key]; ok { //todo return error
			if param.Aggregate != nil && *param.Aggregate {
				aggsMap[param.Key] = struct{}{}
			}
			if len(param.Filter) > 0 {
				filtersMap[param.Key] = f.queryFilter(param)
			}
		} else {
			return errors.New("unknown facet key: " + param.Key)
		}
	}
	for _, f := range fs {
		if _, ok := aggsMap[f.Key()]; ok {
			var aggFilters []osquery.Mappable
			for key, filter := range filtersMap {
				if key != f.Key() {
					aggFilters = append(aggFilters, filter)
				}
			}
			req.Aggs(
				osquery.FilterAgg(f.Key(), osquery.Bool().Must(aggFilters...)).
					Aggs(allFacetsMap[f.Key()].queryAgg()),
			)
		}
		if filter, ok := filtersMap[f.Key()]; ok {
			if filter != nil {
				req.PostFilter(filter)
			}
		}
	}
	return nil
}

func (fs facets) res(resultAggs map[string]json.RawMessage) ([]*search.FacetResult, error) {
	var facetResults []*search.FacetResult
	for _, f := range fs {
		rawAgg, ok := resultAggs[f.Key()]
		if !ok {
			continue
		}
		items, err := f.resultAggItems(rawAgg)
		if err != nil {
			return nil, err
		}
		facetResults = append(facetResults, &search.FacetResult{
			Key:   f.Key(),
			Items: items,
		})
	}
	return facetResults, nil
}

var allFacetsMap = func() map[string]facet {
	m := make(map[string]facet)
	for _, f := range allFacets {
		m[f.Key()] = f
	}
	return m
}()

var (
	allFacets = facets{
		contentTypeFacet,
		fileTypeFacet,
		languageFacet,
		torrentSourceFacet,
	}

	contentTypeFacet = termsFacet{
		key:          "content_type",
		valuePath:    "contentType",
		missingValue: "null",
	}

	fileTypeFacet = nestedFacet{
		key:        "file_type",
		nestedPath: "torrent.files",
		valuePath:  "torrent.files.fileType",
	}

	languageFacet = termsFacet{
		key:       "language",
		valuePath: "languages",
	}

	torrentSourceFacet = nestedFacet{
		key:        "torrent_source",
		nestedPath: "torrent.sources",
		valuePath:  "torrent.sources.source",
	}
)

type nestedFacet struct {
	key        string
	nestedPath string
	valuePath  string
}

func (f nestedFacet) Key() string {
	return f.key
}

func (f nestedFacet) queryAgg() osquery.Aggregation {
	return osquery.NestedAgg("agg", f.nestedPath).Aggs(
		osquery.TermsAgg("value", f.valuePath).Aggs(
			customAgg{
				name: "back_to_root",
				body: map[string]any{
					"reverse_nested": map[string]any{},
				},
			},
		),
	)
}

func (f nestedFacet) queryFilter(facet *search.FacetParam) osquery.Mappable {
	return osquery.Bool().Filter(
		osquery.Nested(
			f.nestedPath,
			osquery.Terms(f.valuePath, func() []any {
				values := make([]any, 0, len(facet.Filter))
				for _, v := range facet.Filter {
					values = append(values, v)
				}
				return values
			}()...),
		),
	)
}

func (nestedFacet) resultAggItems(raw json.RawMessage) ([]*search.FacetResultItem, error) {
	return unmarshalAgg[nestedAgg](raw)
}

type termsFacet struct {
	key          string
	valuePath    string
	missingValue string
}

func (f termsFacet) Key() string {
	return f.key
}

func (termsFacet) Name() string {
	return "agg"
}

func (f termsFacet) Map() map[string]any {
	terms := map[string]any{
		"field": f.valuePath,
	}
	if f.missingValue != "" {
		terms["missing"] = f.missingValue
	}
	return map[string]any{
		"terms": terms,
	}
}

func (f termsFacet) queryAgg() osquery.Aggregation {
	return f
}

func (f termsFacet) queryFilter(facet *search.FacetParam) osquery.Mappable {
	return osquery.Terms(f.valuePath, func() []any {
		values := make([]any, 0, len(facet.Filter))
		for _, v := range facet.Filter {
			values = append(values, v)
		}
		return values
	}()...)
}

func (termsFacet) resultAggItems(raw json.RawMessage) ([]*search.FacetResultItem, error) {
	return unmarshalAgg[responseTermAgg](raw)
}

func unmarshalAgg[T resultAgg](raw json.RawMessage) ([]*search.FacetResultItem, error) {
	var agg T
	err := json.Unmarshal(raw, &agg)
	if err != nil {
		return nil, err
	}
	return agg.resultItems()
}
