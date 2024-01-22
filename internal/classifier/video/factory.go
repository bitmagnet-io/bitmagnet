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
	Resolver lazy.Lazy[classifier.SubClassifier] `group:"content_resolvers"`
}

func New(p Params) Result {
	return Result{
		Resolver: lazy.New(func() (classifier.SubClassifier, error) {
			tmdbClient, err := p.TmdbClient.Get()
			if err != nil {
				return nil, err
			}
			return videoClassifier{
				//config:     classifier.SubClassifierConfig{Key: "video", Priority: 1},
				tmdbClient: tmdbClient,
			}, nil
		}),
	}
}
