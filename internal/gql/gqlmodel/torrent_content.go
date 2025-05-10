package gqlmodel

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
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
	sources := make([]TorrentSourceInfo, 0, len(t.Sources))

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
		options = append(options, torrentContentFacetsOption(*input.Facets))
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

func torrentContentFacetsOption(input gen.TorrentContentFacetsInput) q.Option {
	var qFacets []q.Facet
	if contentType, ok := input.ContentType.ValueOK(); ok {
		qFacets = append(qFacets, torrentContentTypeFacet(*contentType))
	}

	if torrentSource, ok := input.TorrentSource.ValueOK(); ok {
		qFacets = append(qFacets, torrentSourceFacet(*torrentSource))
	}

	if torrentTag, ok := input.TorrentTag.ValueOK(); ok {
		qFacets = append(qFacets, torrentTagFacet(*torrentTag))
	}

	if torrentFileType, ok := input.TorrentFileType.ValueOK(); ok {
		qFacets = append(qFacets, torrentFileTypeFacet(*torrentFileType))
	}

	if language, ok := input.Language.ValueOK(); ok {
		qFacets = append(qFacets, languageFacet(*language))
	}

	if genre, ok := input.Genre.ValueOK(); ok {
		qFacets = append(qFacets, genreFacet(*genre))
	}

	if releaseYear, ok := input.ReleaseYear.ValueOK(); ok {
		qFacets = append(qFacets, releaseYearFacet(*releaseYear))
	}

	if videoResolution, ok := input.VideoResolution.ValueOK(); ok {
		qFacets = append(qFacets, videoResolutionFacet(*videoResolution))
	}

	if videoSource, ok := input.VideoSource.ValueOK(); ok {
		qFacets = append(qFacets, videoSourceFacet(*videoSource))
	}

	return q.WithFacet(qFacets...)
}

func transformTorrentContentSearchResult(
	result q.GenericResult[search.TorrentContentResultItem],
) (TorrentContentSearchResult, error) {
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
	var (
		result gen.TorrentContentAggregations
		err    error
	)

	result.ContentType, err = contentTypeAggs(aggs[search.TorrentContentTypeFacetKey].Items)
	if err != nil {
		return result, err
	}

	result.TorrentSource, err = torrentSourceAggs(aggs[search.TorrentSourceFacetKey].Items)
	if err != nil {
		return result, err
	}

	result.TorrentTag, err = torrentTagAggs(aggs[search.TorrentTagFacetKey].Items)
	if err != nil {
		return result, err
	}

	result.TorrentFileType, err = torrentFileTypeAggs(aggs[search.TorrentFileTypeFacetKey].Items)
	if err != nil {
		return result, err
	}

	result.Language, err = languageAggs(aggs[search.LanguageFacetKey].Items)
	if err != nil {
		return result, err
	}

	result.Genre, err = genreAggs(aggs[search.ContentGenreFacetKey].Items)
	if err != nil {
		return result, err
	}

	result.ReleaseYear, err = releaseYearAggs(aggs[search.ReleaseYearFacetKey].Items)
	if err != nil {
		return result, err
	}

	result.VideoResolution, err = videoResolutionAggs(aggs[search.VideoResolutionFacetKey].Items)
	if err != nil {
		return result, err
	}

	result.VideoSource, err = videoSourceAggs(aggs[search.VideoSourceFacetKey].Items)
	if err != nil {
		return result, err
	}

	return result, nil
}
