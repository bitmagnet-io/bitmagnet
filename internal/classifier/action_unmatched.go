package classifier

import "github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"

const unmatchedName = "unmatched"

type unmatchedAction struct{}

func (unmatchedAction) name() string {
	return unmatchedName
}

var unmatchedPayloadSpec = payloadLiteral[string]{
	literal:     unmatchedName,
	description: "Return a unmatched error for the current torrent",
}

func (unmatchedAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := unmatchedPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}

	path := ctx.path

	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			return ctx.result, classification.RuntimeError{Cause: classification.ErrUnmatched, Path: path}
		},
	}, nil
}

func (unmatchedAction) JSONSchema() JSONSchema {
	return unmatchedPayloadSpec.JSONSchema()
}
