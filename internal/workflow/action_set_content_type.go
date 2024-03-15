package workflow

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/classification"
)

const setContentTypeName = "set_content_type"

type setContentTypeAction struct{}

func (setContentTypeAction) Name() string {
	return setContentTypeName
}

var setContentTypePayloadSpec = payloadSingleKeyValue[model.NullContentType]{
	setContentTypeName,
	payloadMustSucceed[model.NullContentType]{
		payload: contentTypePayloadSpec,
	},
}

func (setContentTypeAction) compileAction(ctx compilerContext) (action, error) {
	contentType, err := setContentTypePayloadSpec.Unmarshal(ctx)
	if err != nil {
		return action{}, ctx.error(err)
	}
	return action{
		func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			cl.ContentType = contentType
			return cl, nil
		},
	}, nil
}
