package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const setContentTypeName = "set_content_type"

type setContentTypeAction struct{}

func (setContentTypeAction) name() string {
	return setContentTypeName
}

var setContentTypeSpec = json_spec.SingleKeyValue[model.NullContentType]{
	Key: setContentTypeName,
	ValueSpec: json_spec.MustSucceed[model.NullContentType]{
		Typed: contentTypeSpec,
	},
	Description: "Set the content type of the current torrent",
}

func (setContentTypeAction) compile(ctx compilerContext) (action, error) {
	contentType, err := setContentTypeSpec.Parse(ctx.jsonSpec)
	if err != nil {
		return action{}, ctx.Error(err)
	}

	return action{
		func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			cl.ContentType = contentType
			return cl, nil
		},
	}, nil
}

func (setContentTypeAction) JSONSchema() json_schema.JSONSchema {
	return setContentTypeSpec.JSONSchema()
}
