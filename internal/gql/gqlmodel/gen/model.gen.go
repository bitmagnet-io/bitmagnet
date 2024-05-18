// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gen

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type ContentTypeAgg struct {
	Value      *model.ContentType `json:"value,omitempty"`
	Label      string             `json:"label"`
	Count      int                `json:"count"`
	IsEstimate bool               `json:"isEstimate"`
}

type ContentTypeFacetInput struct {
	Aggregate graphql.Omittable[*bool]                `json:"aggregate,omitempty"`
	Filter    graphql.Omittable[[]*model.ContentType] `json:"filter,omitempty"`
}

type GenreAgg struct {
	Value      string `json:"value"`
	Label      string `json:"label"`
	Count      int    `json:"count"`
	IsEstimate bool   `json:"isEstimate"`
}

type GenreFacetInput struct {
	Aggregate graphql.Omittable[*bool]             `json:"aggregate,omitempty"`
	Logic     graphql.Omittable[*model.FacetLogic] `json:"logic,omitempty"`
	Filter    graphql.Omittable[[]string]          `json:"filter,omitempty"`
}

type LanguageAgg struct {
	Value      model.Language `json:"value"`
	Label      string         `json:"label"`
	Count      int            `json:"count"`
	IsEstimate bool           `json:"isEstimate"`
}

type LanguageFacetInput struct {
	Aggregate graphql.Omittable[*bool]            `json:"aggregate,omitempty"`
	Filter    graphql.Omittable[[]model.Language] `json:"filter,omitempty"`
}

type Mutation struct {
}

type Query struct {
}

type ReleaseYearAgg struct {
	Value      *model.Year `json:"value,omitempty"`
	Label      string      `json:"label"`
	Count      int         `json:"count"`
	IsEstimate bool        `json:"isEstimate"`
}

type ReleaseYearFacetInput struct {
	Aggregate graphql.Omittable[*bool]         `json:"aggregate,omitempty"`
	Filter    graphql.Omittable[[]*model.Year] `json:"filter,omitempty"`
}

type SuggestTagsQueryInput struct {
	Prefix     graphql.Omittable[*string]  `json:"prefix,omitempty"`
	Exclusions graphql.Omittable[[]string] `json:"exclusions,omitempty"`
}

type SystemQuery struct {
	Version string `json:"version"`
}

type TorrentContentAggregations struct {
	ContentType     []ContentTypeAgg     `json:"contentType,omitempty"`
	TorrentSource   []TorrentSourceAgg   `json:"torrentSource,omitempty"`
	TorrentTag      []TorrentTagAgg      `json:"torrentTag,omitempty"`
	TorrentFileType []TorrentFileTypeAgg `json:"torrentFileType,omitempty"`
	Language        []LanguageAgg        `json:"language,omitempty"`
	Genre           []GenreAgg           `json:"genre,omitempty"`
	ReleaseYear     []ReleaseYearAgg     `json:"releaseYear,omitempty"`
	VideoResolution []VideoResolutionAgg `json:"videoResolution,omitempty"`
	VideoSource     []VideoSourceAgg     `json:"videoSource,omitempty"`
}

type TorrentContentFacetsInput struct {
	ContentType     graphql.Omittable[*ContentTypeFacetInput]     `json:"contentType,omitempty"`
	TorrentSource   graphql.Omittable[*TorrentSourceFacetInput]   `json:"torrentSource,omitempty"`
	TorrentTag      graphql.Omittable[*TorrentTagFacetInput]      `json:"torrentTag,omitempty"`
	TorrentFileType graphql.Omittable[*TorrentFileTypeFacetInput] `json:"torrentFileType,omitempty"`
	Language        graphql.Omittable[*LanguageFacetInput]        `json:"language,omitempty"`
	Genre           graphql.Omittable[*GenreFacetInput]           `json:"genre,omitempty"`
	ReleaseYear     graphql.Omittable[*ReleaseYearFacetInput]     `json:"releaseYear,omitempty"`
	VideoResolution graphql.Omittable[*VideoResolutionFacetInput] `json:"videoResolution,omitempty"`
	VideoSource     graphql.Omittable[*VideoSourceFacetInput]     `json:"videoSource,omitempty"`
}

type TorrentContentOrderByInput struct {
	Field      TorrentContentOrderBy    `json:"field"`
	Descending graphql.Omittable[*bool] `json:"descending,omitempty"`
}

type TorrentFileTypeAgg struct {
	Value      model.FileType `json:"value"`
	Label      string         `json:"label"`
	Count      int            `json:"count"`
	IsEstimate bool           `json:"isEstimate"`
}

type TorrentFileTypeFacetInput struct {
	Aggregate graphql.Omittable[*bool]             `json:"aggregate,omitempty"`
	Logic     graphql.Omittable[*model.FacetLogic] `json:"logic,omitempty"`
	Filter    graphql.Omittable[[]model.FileType]  `json:"filter,omitempty"`
}

type TorrentSourceAgg struct {
	Value      string `json:"value"`
	Label      string `json:"label"`
	Count      int    `json:"count"`
	IsEstimate bool   `json:"isEstimate"`
}

type TorrentSourceFacetInput struct {
	Aggregate graphql.Omittable[*bool]             `json:"aggregate,omitempty"`
	Logic     graphql.Omittable[*model.FacetLogic] `json:"logic,omitempty"`
	Filter    graphql.Omittable[[]string]          `json:"filter,omitempty"`
}

type TorrentTagAgg struct {
	Value      string `json:"value"`
	Label      string `json:"label"`
	Count      int    `json:"count"`
	IsEstimate bool   `json:"isEstimate"`
}

type TorrentTagFacetInput struct {
	Aggregate graphql.Omittable[*bool]             `json:"aggregate,omitempty"`
	Logic     graphql.Omittable[*model.FacetLogic] `json:"logic,omitempty"`
	Filter    graphql.Omittable[[]string]          `json:"filter,omitempty"`
}

type VideoResolutionAgg struct {
	Value      *model.VideoResolution `json:"value,omitempty"`
	Label      string                 `json:"label"`
	Count      int                    `json:"count"`
	IsEstimate bool                   `json:"isEstimate"`
}

type VideoResolutionFacetInput struct {
	Aggregate graphql.Omittable[*bool]                    `json:"aggregate,omitempty"`
	Filter    graphql.Omittable[[]*model.VideoResolution] `json:"filter,omitempty"`
}

type VideoSourceAgg struct {
	Value      *model.VideoSource `json:"value,omitempty"`
	Label      string             `json:"label"`
	Count      int                `json:"count"`
	IsEstimate bool               `json:"isEstimate"`
}

type VideoSourceFacetInput struct {
	Aggregate graphql.Omittable[*bool]                `json:"aggregate,omitempty"`
	Filter    graphql.Omittable[[]*model.VideoSource] `json:"filter,omitempty"`
}

type TorrentContentOrderBy string

const (
	TorrentContentOrderByRelevance   TorrentContentOrderBy = "Relevance"
	TorrentContentOrderByPublishedAt TorrentContentOrderBy = "PublishedAt"
	TorrentContentOrderByUpdatedAt   TorrentContentOrderBy = "UpdatedAt"
	TorrentContentOrderBySize        TorrentContentOrderBy = "Size"
	TorrentContentOrderByFiles       TorrentContentOrderBy = "Files"
	TorrentContentOrderBySeeders     TorrentContentOrderBy = "Seeders"
	TorrentContentOrderByLeechers    TorrentContentOrderBy = "Leechers"
	TorrentContentOrderByName        TorrentContentOrderBy = "Name"
	TorrentContentOrderByInfoHash    TorrentContentOrderBy = "InfoHash"
)

var AllTorrentContentOrderBy = []TorrentContentOrderBy{
	TorrentContentOrderByRelevance,
	TorrentContentOrderByPublishedAt,
	TorrentContentOrderByUpdatedAt,
	TorrentContentOrderBySize,
	TorrentContentOrderByFiles,
	TorrentContentOrderBySeeders,
	TorrentContentOrderByLeechers,
	TorrentContentOrderByName,
	TorrentContentOrderByInfoHash,
}

func (e TorrentContentOrderBy) IsValid() bool {
	switch e {
	case TorrentContentOrderByRelevance, TorrentContentOrderByPublishedAt, TorrentContentOrderByUpdatedAt, TorrentContentOrderBySize, TorrentContentOrderByFiles, TorrentContentOrderBySeeders, TorrentContentOrderByLeechers, TorrentContentOrderByName, TorrentContentOrderByInfoHash:
		return true
	}
	return false
}

func (e TorrentContentOrderBy) String() string {
	return string(e)
}

func (e *TorrentContentOrderBy) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TorrentContentOrderBy(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TorrentContentOrderBy", str)
	}
	return nil
}

func (e TorrentContentOrderBy) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
