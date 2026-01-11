package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

/*
ENUM(

	release_year,
	content_genre,
	language,
	content_type,
	video_3d,
	video_codec,
	video_modifier,
	video_resolution,
	video_source,
	file_type,
	torrent_source,
	tag,

)
*/
type Facet string

var facetLabels = map[Facet]string{
	"release_year":     "Release Year",
	"content_genre":    "Content Genre",
	"language":         "Language",
	"content_type":     "Content Type",
	"video_3d":         "Video 3D",
	"video_codec":      "Video Codec",
	"video_modifier":   "Video Modifier",
	"video_resolution": "Video Resolution",
	"video_source":     "Video Source",
	"file_type":        "File Type",
	"torrent_source":   "Torrent Source",
	"tag":              "Tag",
}

func (f Facet) Label() string {
	return facetLabels[f]
}

type FacetResultItem struct {
	Value      string
	Label      string
	Count      uint
	IsEstimate bool
}

type FacetResult struct {
	Key   Facet
	Logic model.FacetLogic
	Items []FacetResultItem
}

type FacetParam struct {
	Key       Facet
	Filter    []string
	Aggregate bool
	Logic     model.NullFacetLogic
}
