package adapter

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Search search.Search
}

type Result struct {
	fx.Out
	Client torznab.Client
}

func New(p Params) Result {
	return Result{
		Client: adapter{
			title:        "bitmagnet",
			maxLimit:     100,
			defaultLimit: 100,
			search:       p.Search,
		},
	}
}

type adapter struct {
	title        string
	maxLimit     uint
	defaultLimit uint
	search       search.Search
}
