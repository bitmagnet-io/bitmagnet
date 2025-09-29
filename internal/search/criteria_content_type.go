package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const nameContentType = "contentType"

type CriteriaContentType []model.ContentType

func (CriteriaContentType) criteria() {}

var criteriaSpecContentType = json_spec.SingleKeyValue[[]model.ContentType]{
	Key: nameContentType,
	ValueSpec: json_spec.MustSucceed[[]model.ContentType]{
		Typed: json_spec.List[model.ContentType]{
			ItemSpec: json_spec.Enum[model.ContentType]{
				Values: model.ContentTypeValues(),
			},
		},
	},
}

type definitionContentType struct{}

func (definitionContentType) name() string {
	return nameContentType
}

func (definitionContentType) compile(ctx compilerContext) (Criteria, error) {
	contentTypes, err := criteriaSpecContentType.Parse(ctx.jsonSpec)
	if err != nil {
		return nil, ctx.Error(err)
	}

	return CriteriaContentType(contentTypes), nil
}

func (definitionContentType) resultTypes() map[ResultType]struct{} {
	return map[ResultType]struct{}{
		ResultTypeContent:        {},
		ResultTypeTorrentContent: {},
	}
}

func (definitionContentType) JSONSchema() json_schema.JSONSchema {
	return criteriaSpecContentType.JSONSchema()
}
