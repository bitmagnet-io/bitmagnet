package search

import (
	"encoding/json"

	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const nameContentRef = "contentRef"

type CriteriaContentRef []model.ContentRef

func (c CriteriaContentRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]model.ContentRef{
		nameContentRef: c,
	})
}

func (CriteriaContentRef) criteria() {}

type contentRef struct {
	Type   any
	Source any
	ID     any
}

var specContentRef = json_spec.List[model.ContentRef]{
	ItemSpec: json_spec.Transformer[contentRef, model.ContentRef]{
		Typed: json_spec.Struct[contentRef]{},
		Transform: func(cr contentRef, ctx json_spec.ParseContext) (model.ContentRef, error) {
			ct, err := specContentType.Parse(ctx.Child("type", cr.Type))
			if err != nil {
				return model.ContentRef{}, err
			}
			source, err := stringSpec.Parse(ctx.Child("source", cr.Source))
			if err != nil {
				return model.ContentRef{}, err
			}
			id, err := stringSpec.Parse(ctx.Child("id", cr.ID))
			if err != nil {
				return model.ContentRef{}, err
			}
			return model.ContentRef{
				Type:   ct,
				Source: source,
				ID:     id,
			}, nil
		},
	},
}

var criteriaSpecContentRef = json_spec.SingleKeyValue[[]model.ContentRef]{
	Key:       nameContentRef,
	ValueSpec: json_spec.MustSucceed[[]model.ContentRef]{Typed: specContentRef},
}

type definitionContentRef struct{}

func (definitionContentRef) name() string {
	return nameContentRef
}

func (definitionContentRef) compile(ctx compilerContext) (Criteria, error) {
	contentRefs, err := criteriaSpecContentRef.Parse(ctx.jsonSpec)
	if err != nil {
		return nil, ctx.Error(err)
	}

	return CriteriaContentRef(contentRefs), nil
}

func (definitionContentRef) resultTypes() map[ResultType]struct{} {
	return map[ResultType]struct{}{
		ResultTypeContent:        {},
		ResultTypeTorrentContent: {},
	}
}

func (definitionContentRef) JSONSchema() json_schema.JSONSchema {
	return criteriaSpecContentRef.JSONSchema()
}
