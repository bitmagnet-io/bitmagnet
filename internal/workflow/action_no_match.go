package workflow

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
)

const noMatchName = "no_match"

type noMatchAction struct{}

func (noMatchAction) Name() string {
	return noMatchName
}

var noMatchPayloadSpec = payloadLiteral[string]{noMatchName}

func (noMatchAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := noMatchPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	path := ctx.path
	return action{
		run: func(ctx executionContext) (classifier.Classification, error) {
			return ctx.result, runtimeError{cause: ErrNoMatch, path: path}
		},
	}, nil
}
