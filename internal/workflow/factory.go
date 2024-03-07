package workflow

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
)

type Params struct {
	fx.In
	Search lazy.Lazy[search.Search]
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
		return compiler{
			celEnvOption,
			conditions(
				andCondition{},
				//fileExtensionRatioCondition{},
				//fileTypeStatsCondition{},
				//hasContentTypeCondition{},
				//hasFilesStatusCondition{},
				//includesKeywords{},
				orCondition{},
				expressionCondition{},
			),
			actions(
				attachLocalContentByIdAction{
					search: s,
				},
				deleteAction{},
				findMatchAction{},
				ifElseAction{},
				noMatchAction{},
				noopAction{},
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
