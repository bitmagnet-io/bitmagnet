package workflow

import "github.com/bitmagnet-io/bitmagnet/internal/classifier"

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
		run: func(ctx executionContext) (classifier.Classification, error) {
			return ctx.result, runtimeError{cause: ErrDeleteTorrent, path: path}
		},
	}, nil
}
