package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const nameLangauge = "language"

type CriteriaLanguage []model.Language

func (CriteriaLanguage) criteria() {}

var criteriaSpecLanguage = json_spec.SingleKeyValue[[]model.Language]{
	Key: nameLangauge,
	ValueSpec: json_spec.MustSucceed[[]model.Language]{
		Typed: json_spec.List[model.Language]{
			ItemSpec: json_spec.Enum[model.Language]{
				Values: model.LanguageValues(),
			},
		},
	},
}

type definitionLanguage struct{}

func (definitionLanguage) name() string {
	return nameLangauge
}

func (definitionLanguage) compile(ctx compilerContext) (Criteria, error) {
	languages, err := criteriaSpecLanguage.Parse(ctx.jsonSpec)
	if err != nil {
		return nil, ctx.Error(err)
	}

	return CriteriaLanguage(languages), nil
}

func (definitionLanguage) resultTypes() map[ResultType]struct{} {
	return map[ResultType]struct{}{
		ResultTypeContent:        {},
		ResultTypeTorrentContent: {},
	}
}

func (definitionLanguage) JSONSchema() json_schema.JSONSchema {
	return criteriaSpecLanguage.JSONSchema()
}
