package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const addTagName = "add_tag"

type addTagAction struct{}

func (addTagAction) name() string {
	return addTagName
}

var tagPayloadSpec = payloadTransformer[string, string]{
	spec: payloadGeneric[string]{
		jsonSchema: JSONSchema{
			"type": "string",
		},
	},
	transform: func(str string, _ compilerContext) (string, error) {
		if err := model.ValidateTagName(str); err != nil {
			return "", err
		}
		return str, nil
	},
}

var addTagPayloadSpec = payloadSingleKeyValue[[]string]{
	key: addTagName,
	valueSpec: payloadMustSucceed[[]string]{
		payloadList[string]{
			itemSpec: tagPayloadSpec,
		},
	},
	description: "Add one or more tags to the current torrent",
}

func (addTagAction) compileAction(ctx compilerContext) (action, error) {
	tags, err := addTagPayloadSpec.Unmarshal(ctx)
	if err != nil {
		return action{}, ctx.error(err)
	}
	return action{
		func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			if cl.Tags == nil {
				cl.Tags = make(map[string]struct{})
			}
			for _, tag := range tags {
				cl.Tags[tag] = struct{}{}
			}
			return cl, nil
		},
	}, nil
}

func (addTagAction) JSONSchema() JSONSchema {
	return addTagPayloadSpec.JSONSchema()
}
