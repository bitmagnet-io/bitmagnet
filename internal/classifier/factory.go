package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Config     Config
	Search     lazy.Lazy[search.Search]
	TmdbClient lazy.Lazy[tmdb.Client]
}

type Result struct {
	fx.Out
	SourceProvider sourceProvider
	Compiler       lazy.Lazy[Compiler]
	Runner         lazy.Lazy[Runner]
}

func New(params Params) Result {
	sp := newSourceProvider(params.Config)
	lc := lazy.New(func() (Compiler, error) {
		s, err := params.Search.Get()
		if err != nil {
			return nil, err
		}
		tmdbClient, err := params.TmdbClient.Get()
		if err != nil {
			return nil, err
		}
		return compiler{
			celEnvOption,
			conditions(
				andCondition{},
				orCondition{},
				expressionCondition{},
			),
			actions(
				attachLocalContentByIdAction{
					searchAction: searchAction{
						search: s,
					},
				},
				attachLocalContentBySearchAction{
					searchAction: searchAction{
						search: s,
					},
				},
				attachTmdbContentByIdAction{
					tmdbAction: tmdbAction{
						client: tmdbClient,
					},
				},
				attachTmdbContentBySearchAction{
					tmdbAction: tmdbAction{
						client: tmdbClient,
					},
				},
				deleteAction{},
				findMatchAction{},
				ifElseAction{},
				noMatchAction{},
				parseDateAction{},
				parseVideoContentAction{},
				setContentTypeAction{},
			),
		}, nil
	})
	return Result{
		SourceProvider: sp,
		Compiler:       lc,
		Runner: lazy.New(func() (Runner, error) {
			src, err := sp.source()
			if err != nil {
				return nil, err
			}
			c, err := lc.Get()
			if err != nil {
				return nil, err
			}
			return c.Compile(src)
		}),
	}
}
