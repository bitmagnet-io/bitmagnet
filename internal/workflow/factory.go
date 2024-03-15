package workflow

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
)

type Params struct {
	fx.In
	Search     lazy.Lazy[search.Search]
	TmdbClient lazy.Lazy[tmdb.Client]
}

type Result struct {
	fx.Out
	Compiler lazy.Lazy[Compiler]
	Workflow lazy.Lazy[Workflow]
}

func New(params Params) Result {
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
		Compiler: lc,
		Workflow: lazy.New(func() (Workflow, error) {
			c, err := lc.Get()
			if err != nil {
				return nil, err
			}
			rawWorkflow := make(map[string]interface{})
			parseErr := yaml.Unmarshal([]byte(workflowDefaultYaml), &rawWorkflow)
			if parseErr != nil {
				return nil, parseErr
			}
			return c.Compile(rawWorkflow)
		}),
	}
}
