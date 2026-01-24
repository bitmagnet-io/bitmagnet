package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const addTagName = "add_tag"

type addTagAction struct{}

func (addTagAction) name() string {
	return addTagName
}

var tagSpec = json_spec.Transformer[string, string]{
	Typed: json_spec.Generic[string]{
		Schema: json_schema.MustNew(
			json_schema.Typed(json_schema.TypeString),
		),
	},
	Transform: func(str string, _ json_spec.ParseContext) (string, error) {
		if err := model.ValidateTagName(str); err != nil {
			return "", err
		}
		return str, nil
	},
}

var addTagSpec = json_spec.SingleKeyValue[[]string]{
	Key: addTagName,
	ValueSpec: json_spec.MustSucceed[[]string]{
		Typed: json_spec.List[string]{
			ItemSpec: tagSpec,
		},
	},
	Description: "Add one or more tags to the current torrent",
}

func (addTagAction) compile(ctx compilerContext) (action, error) {
	tags, err := addTagSpec.Parse(ctx.jsonSpec)
	if err != nil {
		return action{}, ctx.Error(err)
	}

	return action{
		func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			for _, tag := range tags {
				cl.Tags[tag] = true
			}

			return cl, nil
		},
	}, nil
}

func (addTagAction) JSONSchema() json_schema.JSONSchema {
	return addTagSpec.JSONSchema()
}
