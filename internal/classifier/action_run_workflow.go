package classifier

import (
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const runWorkflowName = "run_workflow"

type runWorkflowAction struct{}

func (runWorkflowAction) name() string {
	return runWorkflowName
}

var runWorkflowSpec = json_spec.SingleKeyValue[[]string]{
	Key: runWorkflowName,
	ValueSpec: json_spec.MustSucceed[[]string]{
		Typed: json_spec.List[string]{
			ItemSpec: json_spec.Generic[string]{
				Schema: json_schema.MustNew(
					json_schema.Typed(json_schema.TypeString),
					json_schema.MinLength(1),
				),
			},
		},
	},
	Description: "Run a different workflow within the current workflow",
}

func (runWorkflowAction) compile(ctx compilerContext) (action, error) {
	names, err := runWorkflowSpec.Parse(ctx.jsonSpec)
	if err != nil {
		return action{}, ctx.Error(err)
	}

	for _, name := range names {
		if _, ok := ctx.workflowNames[Workflow(name)]; !ok {
			return action{}, ctx.Fatal(fmt.Errorf("workflow %s not found", name))
		}
	}

	return action{
		func(ctx executionContext) (classification.Result, error) {
			var err error

			cl := ctx.result
			for _, name := range names {
				cl, err = ctx.workflows[Workflow(name)].run(ctx.withResult(cl))
				if err != nil {
					return cl, err
				}
			}

			return cl, nil
		},
	}, nil
}

func (runWorkflowAction) JSONSchema() json_schema.JSONSchema {
	return runWorkflowSpec.JSONSchema()
}
