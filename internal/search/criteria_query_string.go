package search

import (
	"encoding/json"

	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const nameQueryString = "query"

type CriteriaQueryString string

func (c CriteriaQueryString) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		nameQueryString: string(c),
	})
}

func (CriteriaQueryString) criteria() {}

var criteriaSpecQueryString = json_spec.SingleKeyValue[string]{
	Key: nameQueryString,
	ValueSpec: json_spec.MustSucceed[string]{
		Typed: stringSpec,
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
