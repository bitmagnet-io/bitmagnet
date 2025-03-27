package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func TorrentFileTypeCriteria(fileTypes ...model.FileType) query.Criteria {
	var extensions []string
	for _, fileType := range fileTypes {
		extensions = append(extensions, fileType.Extensions()...)
	}

	return TorrentFileExtensionCriteria(extensions...)
}
