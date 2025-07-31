package classifier

import (
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Config      Config
	Search      search.Search
	TmdbClient  tmdb.Client
	TmdbEnabled tmdb.Enabled
}

type Result struct {
	fx.Out
	Compiler Compiler
	Source   Source
	Runner   Runner
}

func New(params Params) (Result, error) {
	src, err := newSourceProvider(params.Config, params.TmdbEnabled).source()
	if err != nil {
		return Result{}, err
	}

	if _, ok := src.Workflows[params.Config.Workflow]; !ok {
		return Result{}, fmt.Errorf("default workflow '%s' not found", params.Config.Workflow)
	}

	cmp := compiler{
		options: []compilerOption{
			compilerFeatures(defaultFeatures),
			celEnvOption,
		},
		dependencies: dependencies{
			search: localSearchSemaphore{
				search:    localSearch{params.Search},
				semaphore: make(chan struct{}, 1),
			},
			tmdbClient: params.TmdbClient,
		},
	}

	rnr, err := cmp.Compile(src)
	if err != nil {
		return Result{}, err
	}

	rnr = runnerSemaphore{
		runner:    rnr,
		semaphore: make(chan struct{}, params.Config.Concurrency),
	}

	return Result{
		Compiler: cmp,
		Source:   src,
		Runner:   rnr,
	}, nil
}
