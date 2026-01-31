package search

import (
	"encoding/json"

	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const nameGenre = "genre"

type CriteriaGenre []string

func (c CriteriaGenre) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]string{
		nameGenre: c,
	})
}

func (CriteriaGenre) criteria() {}

var criteriaSpecGenre = json_spec.SingleKeyValue[[]string]{
	Key: nameGenre,
	ValueSpec: json_spec.MustSucceed[[]string]{
		Typed: json_spec.List[string]{
			ItemSpec: stringSpec,
		},
	},
}

type definitionGenre struct{}

func (definitionGenre) name() string {
	return nameGenre
}

func (definitionGenre) compile(ctx compilerContext) (Criteria, error) {
	genres, err := criteriaSpecGenre.Parse(ctx.jsonSpec)
	if err != nil {
		return nil, ctx.Error(err)
	}

	return CriteriaGenre(genres), nil
}

func (definitionGenre) resultTypes() map[ResultType]struct{} {
	return map[ResultType]struct{}{
		ResultTypeContent:        {},
		ResultTypeTorrentContent: {},
	}
}

func (definitionGenre) JSONSchema() json_schema.JSONSchema {
	return criteriaSpecGenre.JSONSchema()
}
