package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	adapter "github.com/bitmagnet-io/bitmagnet/internal/search"
)

type Adapter struct {
	Search Search
}

func transformGenericParams(params adapter.Params) query.Option {
	var options []query.Option

	if params.QueryString.Valid {
		options = append(options, query.SearchString(params.QueryString.String))
		if len(params.OrderBy) == 0 {
			options = append(options, query.OrderByQueryStringRank())
		}
	}

	offset := uint(0)

	if params.Limit.Valid {
		options = append(options, query.Limit(params.Limit.Uint))

		if params.Page.Valid && params.Page.Uint > 0 {
			offset += (params.Page.Uint - 1) * params.Limit.Uint
		}
	}

	if params.Offset.Valid {
		offset += params.Offset.Uint
	}

	if offset > 0 {
		options = append(options, query.Offset(offset))
	}

	if params.TotalCount.Valid {
		options = append(options, query.WithTotalCount(params.TotalCount.Bool))
	}

	if params.HasNextPage.Valid {
		options = append(options, query.WithHasNextPage(params.HasNextPage.Bool))
	}

	if params.Cached.Valid {
		if params.Cached.Bool {
			options = append(options, query.Cached())
		} else {
			options = append(options, query.CacheWarm())
		}
	}

	if params.AggregationBudget.Valid {
		options = append(options, query.WithAggregationBudget(params.AggregationBudget.Float64))
	}

	return query.Options(options...)
}

func createFacetOptions(facet adapter.FacetParam) []query.FacetOption {
	var options []query.FacetOption
	if facet.Logic.Valid {
		options = append(options, query.FacetUsesLogic(facet.Logic.FacetLogic))
	}

	if facet.Aggregate {
		options = append(options, query.FacetIsAggregated())
	}

	if len(facet.Filter) > 0 {
		options = append(options, query.FacetHasFilter(valuesToStringKeys(facet.Filter)))
	}

	return options
}

func valuesToStringKeys[T ~string](values []T) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[string(value)] = struct{}{}
	}

	return result
}
