package search

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	adapter "github.com/bitmagnet-io/bitmagnet/internal/search"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

func (a Adapter) TorrentFiles(ctx context.Context, params adapter.Params) (adapter.TorrentFilesResult, error) {
	options := []query.Option{
		transformGenericParams(params),
	}

	if params.Criteria != nil {
		criteria, err := transformTorrentFilesCriteria(params.Criteria)
		if err != nil {
			return adapter.TorrentFilesResult{}, err
		}

		options = append(options, query.Where(criteria))
	}

	if len(params.Facets) > 0 {
		option, err := transformTorrentFilesFacets(params.Facets)
		if err != nil {
			return adapter.TorrentFilesResult{}, err
		}

		options = append(options, option)
	}

	if len(params.OrderBy) > 0 {
		clauses := make([]query.OrderByColumn, 0, len(params.OrderBy))

		for _, ob := range params.OrderBy {
			column, err := adapter.ParseTorrentFilesOrderBy(ob.Key)
			if err != nil {
				return adapter.TorrentFilesResult{}, fmt.Errorf("unknown order by column: %s", ob.Key)
			}

			clauses = append(clauses, TorrentFilesOrderByClauses(column, ob.Descending)...)
		}

		options = append(options, query.OrderBy(clauses...))
	}

	return a.Search.TorrentFiles(ctx, options...)
}

func transformTorrentFilesCriteria(criteria adapter.Criteria) (query.Criteria, error) {
	switch c := criteria.(type) {
	case adapter.And:
		conds, err := slice.MapErr(c, transformTorrentFilesCriteria)
		if err != nil {
			return nil, err
		}

		return query.And(conds...), nil
	case adapter.Or:
		conds, err := slice.MapErr(c, transformTorrentFilesCriteria)
		if err != nil {
			return nil, err
		}

		return query.Or(conds...), nil
	case adapter.Not:
		conds, err := slice.MapErr(c, transformTorrentFilesCriteria)
		if err != nil {
			return nil, err
		}

		return query.Not(conds...), nil
	case adapter.CriteriaInfoHash:
		return TorrentFileInfoHashCriteria(c...), nil
	default:
		return nil, fmt.Errorf("unsupported criteria: %T", criteria)
	}
}

func transformTorrentFilesFacets(facets []adapter.FacetParam) (query.Option, error) {
	qFacets := make([]query.Facet, 0, len(facets))
	seenFacets := make(map[adapter.Facet]struct{})

	for _, facet := range facets {
		if _, ok := seenFacets[facet.Key]; ok {
			return nil, fmt.Errorf("facet key is repeated: %s", facet.Key)
		}

		facet, err := transformTorrentFilesFacet(facet)
		if err != nil {
			return nil, err
		}

		qFacets = append(qFacets, facet)
	}

	return query.WithFacet(qFacets...), nil
}

func transformTorrentFilesFacet(facet adapter.FacetParam) (query.Facet, error) {
	switch facet.Key {
	case adapter.FacetFileType:
		return TorrentFileTypeFacet(createFacetOptions(facet)...), nil
	default:
		return nil, fmt.Errorf("unknown facet: %s", facet.Key)
	}
}
