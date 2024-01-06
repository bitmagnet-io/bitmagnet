package adapter

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Search lazy.Lazy[search.Search]
}

type Result struct {
	fx.Out
	Client lazy.Lazy[torznab.Client]
}

func New(p Params) Result {
	return Result{
		Client: lazy.New[torznab.Client](func() (torznab.Client, error) {
			s, err := p.Search.Get()
			if err != nil {
				return nil, err
			}
			return adapter{
				title:        "bitmagnet",
				maxLimit:     100,
				defaultLimit: 100,
				search:       s,
			}, nil
		}),
	}
}

type adapter struct {
	title        string
	maxLimit     uint
	defaultLimit uint
	search       search.Search
}
