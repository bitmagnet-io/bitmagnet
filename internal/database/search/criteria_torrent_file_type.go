package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

func TorrentFileTypeCriteria(fileTypes ...model.FileType) query.Criteria {
	return TorrentFileExtensionCriteria(slice.FlatMap(fileTypes, func(ft model.FileType) []string {
		return ft.Extensions()
	})...)
}
