package gqlmodel

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	q "github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"time"
)

type TorrentContentQuery struct {
	TorrentContentSearch search.TorrentContentSearch
}

type TorrentContent struct {
	ID              string
	InfoHash        protocol.ID
	ContentType     model.NullContentType
	ContentSource   model.NullString
	ContentID       model.NullString
	Title           string
	Languages       []model.Language `json:"omitempty"`
	Episodes        *Episodes
	VideoResolution model.NullVideoResolution
	VideoSource     model.NullVideoSource
	VideoCodec      model.NullVideoCodec
	Video3D         model.NullVideo3D
	VideoModifier   model.NullVideoModifier
	ReleaseGroup    model.NullString
	SearchString    string
	Seeders         model.NullUint
	Leechers        model.NullUint
	PublishedAt     time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Torrent         model.Torrent
	Content         *model.Content
}

type Episodes struct {
	Label   string
	Seasons []model.Season `json:"omitempty"`
}

func NewTorrentContentFromResultItem(item search.TorrentContentResultItem) TorrentContent {
	c := TorrentContent{
		ID:              item.ID,
		InfoHash:        item.InfoHash,
		ContentType:     item.ContentType,
		ContentSource:   item.ContentSource,
		ContentID:       item.ContentID,
		Title:           item.Title(),
		VideoResolution: item.VideoResolution,
		VideoSource:     item.VideoSource,
		VideoCodec:      item.VideoCodec,
		Video3D:         item.Video3D,
		VideoModifier:   item.VideoModifier,
		ReleaseGroup:    item.ReleaseGroup,
		Seeders:         item.Seeders,
		Leechers:        item.Leechers,
		PublishedAt:     item.PublishedAt,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
		Torrent:         item.Torrent,
	}
	if item.Content.ID != "" {
		c.Content = &item.Content
	}
	languages := item.Languages.Slice()
	if len(languages) > 0 {
		c.Languages = languages
	}
	if len(item.Episodes) > 0 {
		c.Episodes = &Episodes{
			Label:   item.Episodes.String(),
			Seasons: item.Episodes.SeasonEntries(),
		}
	}
	return c
}

type TorrentSourceInfo struct {
	Key      string
	Name     string
	ImportID model.NullString
	Seeders  model.NullUint
	Leechers model.NullUint
}

func TorrentSourceInfosFromTorrent(t model.Torrent) []TorrentSourceInfo {
	var sources []TorrentSourceInfo
	for _, s := range t.Sources {
		sources = append(sources, TorrentSourceInfo{
			Key:      s.Source,
			Name:     s.TorrentSource.Name,
			ImportID: s.ImportID,
			Seeders:  s.Seeders,
			Leechers: s.Leechers,
		})
	}
	return sources
}

type TorrentContentSearchQueryInput struct {
	q.SearchParams
	Facets     *gen.TorrentContentFacetsInput
	OrderBy    []gen.TorrentContentOrderByInput
	InfoHashes graphql.Omittable[[]protocol.ID]
}

type TorrentContentSearchResult struct {
	TotalCount           uint
	TotalCountIsEstimate bool
	HasNextPage          bool
	Items                []TorrentContent
	Aggregations         gen.TorrentContentAggregations
}

func (t TorrentContentQuery) Search(
	ctx context.Context,
	input TorrentContentSearchQueryInput,
) (TorrentContentSearchResult, error) {
	options := []q.Option{
		q.DefaultOption(),
		search.TorrentContentCoreJoins(),
		search.HydrateTorrentContentContent(),
		search.HydrateTorrentContentTorrent(),
	}
	options = append(options, input.Option())
	hasQueryString := input.QueryString.Valid
	if input.Facets != nil {
		var qFacets []q.Facet
		if contentType, ok := input.Facets.ContentType.ValueOK(); ok {
			qFacets = append(qFacets, torrentContentTypeFacet(*contentType))
		}
		if torrentSource, ok := input.Facets.TorrentSource.ValueOK(); ok {
			qFacets = append(qFacets, torrentSourceFacet(*torrentSource))
		}
		if torrentTag, ok := input.Facets.TorrentTag.ValueOK(); ok {
			qFacets = append(qFacets, torrentTagFacet(*torrentTag))
		}
		if torrentFileType, ok := input.Facets.TorrentFileType.ValueOK(); ok {
			qFacets = append(qFacets, torrentFileTypeFacet(*torrentFileType))
		}
		if language, ok := input.Facets.Language.ValueOK(); ok {
			qFacets = append(qFacets, languageFacet(*language))
		}
		if genre, ok := input.Facets.Genre.ValueOK(); ok {
			qFacets = append(qFacets, genreFacet(*genre))
		}
		if releaseYear, ok := input.Facets.ReleaseYear.ValueOK(); ok {
			qFacets = append(qFacets, releaseYearFacet(*releaseYear))
		}
		if videoResolution, ok := input.Facets.VideoResolution.ValueOK(); ok {
			qFacets = append(qFacets, videoResolutionFacet(*videoResolution))
		}
		if videoSource, ok := input.Facets.VideoSource.ValueOK(); ok {
			qFacets = append(qFacets, videoSourceFacet(*videoSource))
		}
		options = append(options, q.WithFacet(qFacets...))
	}
	if infoHashes, ok := input.InfoHashes.ValueOK(); ok {
		options = append(options, q.Where(search.TorrentContentInfoHashCriteria(infoHashes...)))
	}
	fullOrderBy := maps.NewInsertMap[search.TorrentContentOrderBy, search.OrderDirection]()
	for _, ob := range input.OrderBy {
		if ob.Field == gen.TorrentContentOrderByFieldRelevance && !hasQueryString {
			continue
		}
		direction := search.OrderDirectionAscending
		if desc, ok := ob.Descending.ValueOK(); ok && *desc {
			direction = search.OrderDirectionDescending
		}
		field, err := search.ParseTorrentContentOrderBy(ob.Field.String())
		if err != nil {
			return TorrentContentSearchResult{}, err
		}
		fullOrderBy.Set(field, direction)
	}
	options = append(options, search.TorrentContentFullOrderBy(fullOrderBy).Option())
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
		TotalCount:           result.TotalCount,
		TotalCountIsEstimate: result.TotalCountIsEstimate,
		HasNextPage:          result.HasNextPage,
		Items:                items,
		Aggregations:         aggs,
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
