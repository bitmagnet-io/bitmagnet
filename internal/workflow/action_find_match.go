package workflow

import (
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/classification"
)

const findMatchName = "find_match"

type findMatchAction struct{}

func (findMatchAction) Name() string {
	return findMatchName
}

var findMatchActionPayloadSpec = payloadSingleKeyValue[[]any]{
	key: findMatchName,
	valueSpec: payloadMustSucceed[[]any]{payloadList[any]{payloadGeneric[any]{
		jsonSchema: map[string]any{
			"type": "any",
		},
	}}},
}

func (findMatchAction) compileAction(ctx compilerContext) (action, error) {
	payload, err := findMatchActionPayloadSpec.Unmarshal(ctx)
	if err != nil {
		return action{}, ctx.error(err)
	}
	actions := make([]action, len(payload))
	for i, actionPayload := range payload {
		a, err := ctx.compileAction(ctx.child(numericPathPart(i), actionPayload))
		if err != nil {
			return action{}, err
		}
		actions[i] = a
	}
	path := ctx.path
	return action{
		func(ctx executionContext) (classification.Result, error) {
			for _, action := range actions {
				result, err := action.run(ctx)
				if err != nil {
					if errors.Is(err, classification.ErrNoMatch) {
						continue
					}
					return classification.Result{}, classification.RuntimeError{
						Cause: err,
						Path:  path,
					}
				} else {
					return result, nil
				}
			}
			return ctx.result, nil
		},
	}, nil
}
