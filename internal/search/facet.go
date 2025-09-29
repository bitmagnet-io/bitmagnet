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

type FacetResultItem struct {
	Value      string
	Label      string
	Count      uint
	IsEstimate bool
}

type FacetResult struct {
	Key   Facet
	Label string
	Logic model.FacetLogic
	Items []FacetResultItem
}

type FacetParam struct {
	Key       Facet
	Filter    []string
	Aggregate bool
	Logic     model.NullFacetLogic
}
