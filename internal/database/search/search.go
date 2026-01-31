package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
)

type Search interface {
	ContentSearch
	QueueJobSearch
	TorrentSearch
	TorrentContentSearch
	TorrentFilesSearch
}

type search struct {
	daoProvider database.DaoProvider
}

func New(dao database.DaoProvider) Search {
	return &search{
		daoProvider: dao,
	}
}
