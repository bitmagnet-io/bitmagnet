package gqlmodel

import (
	"context"
	"time"

	q "github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
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
	Video3d         model.NullVideo3d
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
		Video3d:         item.Video3d,
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

type TorrentSource struct {
	Key      string
	Name     string
	ImportID model.NullString
	Seeders  model.NullUint
	Leechers model.NullUint
}

func TorrentSourcesFromTorrent(t model.Torrent) []TorrentSource {
	var sources []TorrentSource
	for _, s := range t.Sources {
		sources = append(sources, TorrentSource{
			Key:      s.Source,
			Name:     s.TorrentSource.Name,
			ImportID: s.ImportID,
			Seeders:  s.Seeders,
			Leechers: s.Leechers,
		})
	}
	return sources
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
	query *q.SearchParams,
	facets *gen.TorrentContentFacetsInput,
	orderBy []gen.TorrentContentOrderByInput,
) (TorrentContentSearchResult, error) {
	options := []q.Option{
		search.TorrentContentDefaultOption(),
	}
	hasQueryString := false
	if query != nil {
		options = append(options, query.Option())
		hasQueryString = query.QueryString.Valid
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
	fullOrderBy := maps.NewInsertMap[search.TorrentContentOrderBy, search.OrderDirection]()
	for _, ob := range orderBy {
		if ob.Field == gen.TorrentContentOrderByRelevance && !hasQueryString {
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

type TorrentContentByID struct {
	Search search.Search
}

func (t TorrentContentByID) Id(
	ctx context.Context,
	infoHash protocol.ID,
) (TorrentContentSearchResult, error) {
	options := []q.Option{
		q.Where(
			search.TorrentContentInfoHashCriteria(infoHash),
		),
		search.TorrentContentCoreJoins(),
		search.HydrateTorrentContentContent(),
		search.HydrateTorrentContentTorrent(),
		q.Limit(1),
	}

	result, err := t.Search.TorrentContent(ctx, options...)
	if err != nil {
		return TorrentContentSearchResult{}, err
	}

	return transformTorrentContentSearchResult(result)
}
