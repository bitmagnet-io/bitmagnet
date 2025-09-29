package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
)

const nameQueryString = "query"

type CriteriaQueryString string

func (CriteriaQueryString) criteria() {}

var criteriaSpecQueryString = json_spec.SingleKeyValue[string]{
	Key: nameQueryString,
	ValueSpec: json_spec.MustSucceed[string]{
		Typed: json_spec.Generic[string]{
			Schema: json_schema.MustNew(
				json_schema.Typed(json_schema.TypeString),
			),
		},
	},
}

type definitionQueryString struct{}

func (definitionQueryString) name() string {
	return nameQueryString
}

func (definitionQueryString) compile(ctx compilerContext) (Criteria, error) {
	str, err := criteriaSpecQueryString.Parse(ctx.jsonSpec)
	if err != nil {
		return nil, ctx.Error(err)
	}

	return CriteriaQueryString(str), nil
}

func (definitionQueryString) resultTypes() map[ResultType]struct{} {
	return map[ResultType]struct{}{
		ResultTypeContent:        {},
		ResultTypeTorrentContent: {},
	}
}

func (definitionQueryString) JSONSchema() json_schema.JSONSchema {
	return criteriaSpecQueryString.JSONSchema()
}
