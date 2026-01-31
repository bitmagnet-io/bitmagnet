package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const deleteName = "delete"

type deleteAction struct{}

func (deleteAction) name() string {
	return deleteName
}

var deleteSpec = json_spec.Literal[string]{
	Literal:     deleteName,
	Description: "Delete the current torrent",
}

func (deleteAction) compile(ctx compilerContext) (action, error) {
	if _, err := deleteSpec.Parse(ctx.jsonSpec); err != nil {
		return action{}, ctx.Error(err)
	}

	path := ctx.Path

	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			return ctx.result, classification.RuntimeError{
				Cause: classification.ErrDeleteTorrent,
				Path:  path,
			}
		},
	}, nil
}

func (deleteAction) JSONSchema() json_schema.JSONSchema {
	return deleteSpec.JSONSchema()
}
