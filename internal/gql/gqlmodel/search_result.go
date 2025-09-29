package gqlmodel

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/search"
)

type BaseSearchResult[T any] search.Result[T]

func (r BaseSearchResult[_]) IsSearchResult() {}

func (r BaseSearchResult[_]) GetTotalCount() *int {
	if r.TotalCount.Valid {
		totalCount := int(r.TotalCount.Uint)
		return &totalCount
	}

	return nil
}

func (r BaseSearchResult[_]) GetTotalCountIsEstimate() bool {
	return r.TotalCountIsEstimate
}

func (r BaseSearchResult[_]) GetHasNextPage() *bool {
	if r.HasNextPage.Valid {
		return &r.HasNextPage.Bool
	}

	return nil
}

type (
	TorrentContent struct {
		model.TorrentContent
		Content *model.Content
	}

	TorrentContentSearchResult = BaseSearchResult[TorrentContent]
	TorrentFilesSearchResult   = BaseSearchResult[model.TorrentFile]
)
