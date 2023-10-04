package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"go.uber.org/fx"
)

type Search interface {
	ContentSearch
	TorrentSearch
	TorrentContentSearch
}

type search struct {
	q *dao.Query
}

type Params struct {
	fx.In
	Query *dao.Query
}

type Result struct {
	fx.Out
	Search Search
}

func New(params Params) Result {
	return Result{
		Search: &search{
			q: params.Query,
		},
	}
}
