package classifier

import "github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"

const deleteName = "delete"

type deleteAction struct{}

func (deleteAction) Name() string {
	return deleteName
}

var deletePayloadSpec = payloadLiteral[string]{deleteName}

func (deleteAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := deletePayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	path := ctx.path
	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			return ctx.result, classification.RuntimeError{Cause: classification.ErrDeleteTorrent, Path: path}
		},
	}, nil
}
