package classifier

import "github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"

const noMatchName = "no_match"

type noMatchAction struct{}

func (noMatchAction) name() string {
	return noMatchName
}

var noMatchPayloadSpec = payloadLiteral[string]{noMatchName}

func (noMatchAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := noMatchPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	path := ctx.path
	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			return ctx.result, classification.RuntimeError{Cause: classification.ErrNoMatch, Path: path}
		},
	}, nil
}

func (noMatchAction) JsonSchema() JsonSchema {
	return noMatchPayloadSpec.JsonSchema()
}
