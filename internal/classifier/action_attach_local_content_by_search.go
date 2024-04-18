package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
)

const attachLocalContentBySearchName = "attach_local_content_by_search"

type attachLocalContentBySearchAction struct{}

func (attachLocalContentBySearchAction) name() string {
	return attachLocalContentBySearchName
}

var attachLocalContentBySearchPayloadSpec = payloadLiteral[string]{attachLocalContentBySearchName}

func (a attachLocalContentBySearchAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := attachLocalContentBySearchPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			if !cl.ContentType.Valid || !cl.BaseTitle.Valid {
				return cl, classification.ErrNoMatch
			}
			content, err := ctx.search.ContentBySearch(ctx.Context, cl.ContentType.ContentType, cl.BaseTitle.String, cl.Date.Year)
			if err != nil {
				return cl, err
			}
			cl.AttachContent(&content)
			return cl, nil
		},
	}, nil
}

func (attachLocalContentBySearchAction) JsonSchema() JsonSchema {
	return attachLocalContentBySearchPayloadSpec.JsonSchema()
}
