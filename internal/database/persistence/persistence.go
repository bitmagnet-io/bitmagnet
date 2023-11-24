package persistence

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"go.uber.org/fx"
)

type Persistence interface {
	TorrentPersistence
}

type persistence struct {
	q *dao.Query
}

type Params struct {
	fx.In
	Query *dao.Query
}

type Result struct {
	fx.Out
	Persistence Persistence
}

func New(params Params) (Result, error) {
	return Result{
		Persistence: &persistence{
			q: params.Query,
		},
	}, nil
}
