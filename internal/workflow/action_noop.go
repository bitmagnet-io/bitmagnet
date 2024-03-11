package workflow

import "github.com/bitmagnet-io/bitmagnet/internal/processor/classification"

const noopName = "noop"

type noopAction struct{}

func (noopAction) Name() string {
	return noopName
}

var noopPayloadSpec = payloadLiteral[string]{noopName}

func (noopAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := noopPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			return ctx.result, nil
		},
	}, nil
}
