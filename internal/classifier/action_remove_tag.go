package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const removeTagName = "remove_tag"

type removeTagAction struct{}

func (removeTagAction) name() string {
	return removeTagName
}

var removeTagSpec = json_spec.SingleKeyValue[[]string]{
	Key: removeTagName,
	ValueSpec: json_spec.MustSucceed[[]string]{
		json_spec.List[string]{
			ItemSpec: tagSpec,
		},
	},
	Description: "Remove one or more tags from the current torrent",
}

func (removeTagAction) compile(ctx compilerContext) (action, error) {
	tags, err := removeTagSpec.Parse(ctx.jsonSpec)
	if err != nil {
		return action{}, ctx.Error(err)
	}

	return action{
		func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			for _, tag := range tags {
				cl.Tags[tag] = false
			}
			return cl, nil
		},
	}, nil
}

func (removeTagAction) JSONSchema() json_schema.JSONSchema {
	return removeTagSpec.JSONSchema()
}
