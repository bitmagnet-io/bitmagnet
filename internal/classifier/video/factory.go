package video

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/video/tmdb"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	TmdbClient lazy.Lazy[tmdb.Client]
}

type Result struct {
	fx.Out
	Resolver lazy.Lazy[classifier.SubResolver] `group:"content_resolvers"`
}

func New(p Params) Result {
	return Result{
		Resolver: lazy.New(func() (classifier.SubResolver, error) {
			tmdbClient, err := p.TmdbClient.Get()
			if err != nil {
				return nil, err
			}
			return videoResolver{
				config:     classifier.SubResolverConfig{Key: "video", Priority: 1},
				tmdbClient: tmdbClient,
			}, nil
		}),
	}
}
