package tmdb

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Search     lazy.Lazy[search.Search]
	TmdbClient lazy.Lazy[tmdb.Client]
	Logger     *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Client lazy.Lazy[Client]
}

func New(p Params) Result {
	return Result{
		Client: lazy.New(func() (Client, error) {
			s, err := p.Search.Get()
			if err != nil {
				return nil, err
			}
			c, err := p.TmdbClient.Get()
			if err != nil {
				return nil, err
			}
			return &client{
				c: c,
				s: s,
			}, nil
		}),
	}
}
