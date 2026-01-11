package search

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	adapter "github.com/bitmagnet-io/bitmagnet/internal/search"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

func (a Adapter) TorrentContent(ctx context.Context, params adapter.Params) (adapter.TorrentContentResult, error) {
	options := []query.Option{
		TorrentContentDefaultOption(),
		transformGenericParams(params),
	}

	if params.Criteria != nil {
		criteria, err := transformTorrentContentCriteria(params.Criteria)
		if err != nil {
			return adapter.TorrentContentResult{}, err
		}

		options = append(options, query.Where(criteria))
	}

	if len(params.Facets) > 0 {
		option, err := transformTorrentContentFacets(params.Facets)
		if err != nil {
			return adapter.TorrentContentResult{}, err
		}

		options = append(options, option)
	}

	if len(params.OrderBy) > 0 {
		clauses := make([]query.OrderByColumn, 0, len(params.OrderBy))

		for _, ob := range params.OrderBy {
			column, err := adapter.ParseTorrentContentOrderBy(ob.Key)
			if err != nil {
				return adapter.TorrentContentResult{}, fmt.Errorf("unknown order by column: %s", ob.Key)
			}

			clauses = append(clauses, TorrentContentOrderByClauses(column, ob.Descending)...)
		}

		options = append(options, query.OrderBy(clauses...))
	}

	return a.Search.TorrentContent(ctx, options...)
}

func transformTorrentContentCriteria(criteria adapter.Criteria) (query.Criteria, error) {
	switch c := criteria.(type) {
	case adapter.And:
		conds, err := slice.MapErr(c, transformTorrentContentCriteria)
		if err != nil {
			return nil, err
		}

		return query.And(conds...), nil
	case adapter.Or:
		conds, err := slice.MapErr(c, transformTorrentContentCriteria)
		if err != nil {
			return nil, err
		}

		return query.Or(conds...), nil
	case adapter.Not:
		conds, err := slice.MapErr(c, transformTorrentContentCriteria)
		if err != nil {
			return nil, err
		}

		return query.Not(conds...), nil
	case adapter.CriteriaContentType:
		return TorrentContentTypeCriteria(c...), nil
	case adapter.CriteriaContentRef:
		return query.Or(
			slice.Map(c, func(ref model.ContentRef) query.Criteria {
				// todo: Clean this up
				if ref.Source == "tmdb" {
					return ContentCanonicalIdentifierCriteria(model.ContentRef{
						Type:   ref.Type,
						Source: ref.Source,
						ID:     ref.ID,
					})
				}

				return ContentIdentifierCriteria(model.ContentRef{
					Type:   ref.Type,
					Source: ref.Source,
					ID:     ref.ID,
				})
			})...,
		), nil
	case adapter.CriteriaInfoHash:
		return TorrentContentInfoHashCriteria(c...), nil
	case adapter.CriteriaGenre:
		return query.Or(TorrentContentGenreFacet().Criteria(valuesToStringKeys(c))...), nil
	case adapter.CriteriaLanguage:
		return query.Or(TorrentContentLanguageFacet().Criteria(valuesToStringKeys(c))...), nil
	case adapter.CriteriaTag:
		return query.Or(TorrentTagsFacet().Criteria(valuesToStringKeys(c))...), nil
	default:
		return nil, fmt.Errorf("unsupported criteria: %T", criteria)
	}
}

func transformTorrentContentFacets(facets []adapter.FacetParam) (query.Option, error) {
	qFacets := make([]query.Facet, 0, len(facets))
	seenFacets := make(map[adapter.Facet]struct{})

	for _, facet := range facets {
		if _, ok := seenFacets[facet.Key]; ok {
			return nil, fmt.Errorf("facet key is repeated: %s", facet.Key)
		}
		facet, err := transformTorrentContentFacet(facet)
		if err != nil {
			return nil, err
		}

		qFacets = append(qFacets, facet)
	}

	return query.WithFacet(qFacets...), nil
}

func transformTorrentContentFacet(facet adapter.FacetParam) (query.Facet, error) {
	switch facet.Key {
	case adapter.FacetContentType:
		return TorrentContentTypeFacet(createFacetOptions(facet)...), nil
	case adapter.FacetTorrentSource:
		return TorrentSourceFacet(createFacetOptions(facet)...), nil
	case adapter.FacetTag:
		return TorrentTagsFacet(createFacetOptions(facet)...), nil
	case adapter.FacetFileType:
		return TorrentFileTypeFacet(createFacetOptions(facet)...), nil
	case adapter.FacetLanguage:
		return TorrentContentLanguageFacet(createFacetOptions(facet)...), nil
	case adapter.FacetContentGenre:
		return TorrentContentGenreFacet(createFacetOptions(facet)...), nil
	case adapter.FacetReleaseYear:
		return ReleaseYearFacet(createFacetOptions(facet)...), nil
	case adapter.FacetVideoResolution:
		return VideoResolutionFacet(createFacetOptions(facet)...), nil
	case adapter.FacetVideoSource:
		return VideoSourceFacet(createFacetOptions(facet)...), nil
	default:
		return nil, fmt.Errorf("unknown facet: %s", facet.Key)
	}
}
