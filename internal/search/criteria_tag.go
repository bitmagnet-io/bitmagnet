package search

import (
	"encoding/json"

	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
)

const nameTag = "tag"

type CriteriaTag []string

func (c CriteriaTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]string{
		nameTag: c,
	})
}

func (CriteriaTag) criteria() {}

var criteriaSpecTag = json_spec.SingleKeyValue[[]string]{
	Key: nameTag,
	ValueSpec: json_spec.MustSucceed[[]string]{
		Typed: json_spec.List[string]{
			ItemSpec: json_spec.Generic[string]{
				Schema: json_schema.MustNew(
					json_schema.Typed(json_schema.TypeString),
				),
			},
		},
	},
}

type definitionTag struct{}

func (definitionTag) name() string {
	return nameTag
}

func (definitionTag) compile(ctx compilerContext) (Criteria, error) {
	tags, err := criteriaSpecTag.Parse(ctx.jsonSpec)
	if err != nil {
		return nil, ctx.Error(err)
	}

	return CriteriaTag(tags), nil
}

func (definitionTag) resultTypes() map[ResultType]struct{} {
	return map[ResultType]struct{}{
		ResultTypeTorrentContent: {},
	}
}

func (definitionTag) JSONSchema() json_schema.JSONSchema {
	return criteriaSpecTag.JSONSchema()
}
