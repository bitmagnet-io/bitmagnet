package gqlmodel

import (
	"context"
	q "github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
)

type TorrentContentSearchResult struct {
	TotalCount   uint
	Items        []TorrentContent
	Aggregations gen.TorrentContentAggregations
}

func (t TorrentContentQuery) Search(ctx context.Context, query *q.SearchParams, facets *gen.TorrentContentFacetsInput) (TorrentContentSearchResult, error) {
	options := []q.Option{
		search.TorrentContentDefaultOption(),
	}
	if query != nil {
		options = append(options, query.Option())
	}
	if facets != nil {
		var qFacets []q.Facet
		if contentType, ok := facets.ContentType.ValueOK(); ok {
			qFacets = append(qFacets, torrentContentTypeFacet(*contentType))
		}
		if torrentSource, ok := facets.TorrentSource.ValueOK(); ok {
			qFacets = append(qFacets, torrentSourceFacet(*torrentSource))
		}
		if torrentTag, ok := facets.TorrentTag.ValueOK(); ok {
			qFacets = append(qFacets, torrentTagFacet(*torrentTag))
		}
		if torrentFileType, ok := facets.TorrentFileType.ValueOK(); ok {
			qFacets = append(qFacets, torrentFileTypeFacet(*torrentFileType))
		}
		if language, ok := facets.Language.ValueOK(); ok {
			qFacets = append(qFacets, languageFacet(*language))
		}
		if genre, ok := facets.Genre.ValueOK(); ok {
			qFacets = append(qFacets, genreFacet(*genre))
		}
		if releaseYear, ok := facets.ReleaseYear.ValueOK(); ok {
			qFacets = append(qFacets, releaseYearFacet(*releaseYear))
		}
		if videoResolution, ok := facets.VideoResolution.ValueOK(); ok {
			qFacets = append(qFacets, videoResolutionFacet(*videoResolution))
		}
		if videoSource, ok := facets.VideoSource.ValueOK(); ok {
			qFacets = append(qFacets, videoSourceFacet(*videoSource))
		}
		options = append(options, q.WithFacet(qFacets...))
	}
	result, resultErr := t.TorrentContentSearch.TorrentContent(ctx, options...)
	if resultErr != nil {
		return TorrentContentSearchResult{}, resultErr
	}
	return transformTorrentContentSearchResult(result)
}

func transformTorrentContentSearchResult(result q.GenericResult[search.TorrentContentResultItem]) (TorrentContentSearchResult, error) {
	aggs, aggsErr := transformTorrentContentAggregations(result.Aggregations)
	if aggsErr != nil {
		return TorrentContentSearchResult{}, aggsErr
	}
	items := make([]TorrentContent, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, NewTorrentContentFromResultItem(item))
	}
	return TorrentContentSearchResult{
		TotalCount:   result.TotalCount,
		Items:        items,
		Aggregations: aggs,
	}, nil
}

func transformTorrentContentAggregations(aggs q.Aggregations) (gen.TorrentContentAggregations, error) {
	a := gen.TorrentContentAggregations{}
	if contentTypes, ok := aggs[search.TorrentContentTypeFacetKey]; ok {
		agg, err := contentTypeAggs(contentTypes.Items)
		if err != nil {
			return a, err
		}
		a.ContentType = agg
	}
	if torrentSources, ok := aggs[search.TorrentSourceFacetKey]; ok {
		agg, err := torrentSourceAggs(torrentSources.Items)
		if err != nil {
			return a, err
		}
		a.TorrentSource = agg
	}
	if torrentTags, ok := aggs[search.TorrentTagFacetKey]; ok {
		agg, err := torrentTagAggs(torrentTags.Items)
		if err != nil {
			return a, err
		}
		a.TorrentTag = agg
	}
	if fileTypes, ok := aggs[search.TorrentFileTypeFacetKey]; ok {
		agg, err := torrentFileTypeAggs(fileTypes.Items)
		if err != nil {
			return a, err
		}
		a.TorrentFileType = agg
	}
	if languages, ok := aggs[search.LanguageFacetKey]; ok {
		agg, err := languageAggs(languages.Items)
		if err != nil {
			return a, err
		}
		a.Language = agg
	}
	if genres, ok := aggs[search.ContentGenreFacetKey]; ok {
		agg, err := genreAggs(genres.Items)
		if err != nil {
			return a, err
		}
		a.Genre = agg
	}
	if releaseYears, ok := aggs[search.ReleaseYearFacetKey]; ok {
		agg, err := releaseYearAggs(releaseYears.Items)
		if err != nil {
			return a, err
		}
		a.ReleaseYear = agg
	}
	if videoResolutions, ok := aggs[search.VideoResolutionFacetKey]; ok {
		agg, err := videoResolutionAggs(videoResolutions.Items)
		if err != nil {
			return a, err
		}
		a.VideoResolution = agg
	}
	if videoSource, ok := aggs[search.VideoSourceFacetKey]; ok {
		agg, err := videoSourceAggs(videoSource.Items)
		if err != nil {
			return a, err
		}
		a.VideoSource = agg
	}
	return a, nil
}
