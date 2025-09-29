package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
)

const nameGenre = "genre"

type CriteriaGenre []string

func (CriteriaGenre) criteria() {}

var criteriaSpecGenre = json_spec.SingleKeyValue[[]string]{
	Key: nameGenre,
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
