package workflow

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
)

const attachTmdbContentByIdName = "attach_tmdb_content_by_id"

type attachTmdbContentByIdAction struct {
	client tmdb.Client
}

func (attachTmdbContentByIdAction) Name() string {
	return attachTmdbContentByIdName
}

var attachTmdbContentByIdPayloadSpec = payloadLiteral[string]{attachTmdbContentByIdName}

func (a attachTmdbContentByIdAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := attachTmdbContentByIdPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	return action{
		run: func(ctx executionContext) (classifier.Classification, error) {
			return classifier.Classification{}, nil
		},
	}, nil
}
