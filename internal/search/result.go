package search

import "github.com/bitmagnet-io/bitmagnet/internal/model"

/*
ENUM(

		torrent,
	  torrent_content,
	  torrent_file,
	  content,

)
*/
type ResultType string

type Result[T any] struct {
	TotalCount           model.NullUint
	TotalCountIsEstimate bool
	HasNextPage          model.NullBool
	Items                []T
	Facets               []FacetResult
}

type ContentResult = Result[ContentResultItem]

type ContentResultItem struct {
	QueryStringRank float64
	model.Content
}

type TorrentContentResultItem struct {
	QueryStringRank float64
	model.TorrentContent
}

type TorrentContentResult = Result[TorrentContentResultItem]

type TorrentFilesResult = Result[model.TorrentFile]

type TorrentsResult = Result[model.Torrent]
