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
	"github.com/vektah/gqlparser/v2/ast"
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

func isPathSelectedCtx(ctx *graphql.OperationContext, selSet ast.SelectionSet, path ...string) bool {
	for _, f := range graphql.CollectFields(ctx, selSet, nil) {
		if f.Name == path[0] {
			if len(path) == 1 {
				return true
			}
			return isPathSelectedCtx(ctx, f.Selections, path[1:]...)
		}
	}
	return false
}

func isPathSelected(ctx context.Context, first string, rest ...string) bool {
	for _, f1 := range graphql.CollectFieldsCtx(ctx, nil) {
		if f1.Name == first {
			if len(rest) == 0 {
				return true
			}
			return isPathSelectedCtx(graphql.GetOperationContext(ctx), f1.Selections, rest...)
		}
	}
	return false
}

func (t TorrentContentQuery) Search(
	ctx context.Context,
	query TorrentContentSearchQueryInput,
) (TorrentContentSearchResult, error) {
	options := []q.Option{
		q.DefaultOption(),
		search.TorrentContentCoreJoins(),
		search.HydrateTorrentContentContent(),
		search.HydrateTorrentContentTorrent(),
	}
	options = append(options, query.Option())
	hasQueryString := query.QueryString.Valid
	if query.Facets != nil {
		var qFacets []q.Facet
		if contentType, ok := query.Facets.ContentType.ValueOK(); ok {
			qFacets = append(qFacets, torrentContentTypeFacet(*contentType))
		}
		if torrentSource, ok := query.Facets.TorrentSource.ValueOK(); ok {
			qFacets = append(qFacets, torrentSourceFacet(*torrentSource))
		}
		if torrentTag, ok := query.Facets.TorrentTag.ValueOK(); ok {
			qFacets = append(qFacets, torrentTagFacet(*torrentTag))
		}
		if torrentFileType, ok := query.Facets.TorrentFileType.ValueOK(); ok {
			qFacets = append(qFacets, torrentFileTypeFacet(*torrentFileType))
		}
		if language, ok := query.Facets.Language.ValueOK(); ok {
			qFacets = append(qFacets, languageFacet(*language))
		}
		if genre, ok := query.Facets.Genre.ValueOK(); ok {
			qFacets = append(qFacets, genreFacet(*genre))
		}
		if releaseYear, ok := query.Facets.ReleaseYear.ValueOK(); ok {
			qFacets = append(qFacets, releaseYearFacet(*releaseYear))
		}
		if videoResolution, ok := query.Facets.VideoResolution.ValueOK(); ok {
			qFacets = append(qFacets, videoResolutionFacet(*videoResolution))
		}
		if videoSource, ok := query.Facets.VideoSource.ValueOK(); ok {
			qFacets = append(qFacets, videoSourceFacet(*videoSource))
		}
		options = append(options, q.WithFacet(qFacets...))
	}
	if infoHashes, ok := query.InfoHashes.ValueOK(); ok {
		options = append(options, q.Where(search.TorrentContentInfoHashCriteria(infoHashes...)))
	}
	fullOrderBy := maps.NewInsertMap[search.TorrentContentOrderBy, search.OrderDirection]()
	for _, ob := range query.OrderBy {
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
