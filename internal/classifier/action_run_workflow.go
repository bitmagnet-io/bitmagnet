package classifier

import (
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
)

const runWorkflowName = "run_workflow"

type runWorkflowAction struct{}

func (runWorkflowAction) name() string {
	return runWorkflowName
}

var runWorkflowPayloadSpec = payloadSingleKeyValue[[]string]{
	key: runWorkflowName,
	valueSpec: payloadMustSucceed[[]string]{
		payloadList[string]{
			itemSpec: payloadGeneric[string]{
				jsonSchema: map[string]interface{}{
					"type":      "string",
					"minLength": 1,
				},
			},
		},
	},
	description: "Run a different workflow within the current workflow",
}

func (runWorkflowAction) compileAction(ctx compilerContext) (action, error) {
	names, err := runWorkflowPayloadSpec.Unmarshal(ctx)
	if err != nil {
		return action{}, ctx.error(err)
	}

	for _, name := range names {
		if _, ok := ctx.workflowNames[name]; !ok {
			return action{}, ctx.fatal(fmt.Errorf("workflow %s not found", name))
		}
	}

	return action{
		func(ctx executionContext) (classification.Result, error) {
			var err error
			cl := ctx.result
			for _, name := range names {
				cl, err = ctx.workflows[name].run(ctx.withResult(cl))
				if err != nil {
					return cl, err
				}
			}
			return cl, nil
		},
	}, nil
}

func (runWorkflowAction) JSONSchema() JSONSchema {
	return runWorkflowPayloadSpec.JSONSchema()
}
