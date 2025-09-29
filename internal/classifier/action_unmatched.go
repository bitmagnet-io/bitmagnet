package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
)

const unmatchedName = "unmatched"

type unmatchedAction struct{}

func (unmatchedAction) name() string {
	return unmatchedName
}

var unmatchedSpec = json_spec.Literal[string]{
	Literal:     unmatchedName,
	Description: "Return a unmatched error for the current torrent",
}

func (unmatchedAction) compile(ctx compilerContext) (action, error) {
	if _, err := unmatchedSpec.Parse(ctx.jsonSpec); err != nil {
		return action{}, ctx.Error(err)
	}

	path := ctx.Path

	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			return ctx.result, classification.RuntimeError{Cause: classification.ErrUnmatched, Path: path}
		},
	}, nil
}

func (unmatchedAction) JSONSchema() json_schema.JSONSchema {
	return unmatchedSpec.JSONSchema()
}
