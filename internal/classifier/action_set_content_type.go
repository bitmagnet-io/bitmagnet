package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const setContentTypeName = "set_content_type"

type setContentTypeAction struct{}

func (setContentTypeAction) name() string {
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

func (setContentTypeAction) JsonSchema() JsonSchema {
	return setContentTypePayloadSpec.JsonSchema()
}
