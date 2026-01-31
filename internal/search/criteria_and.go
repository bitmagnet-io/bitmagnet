package search

import (
	"encoding/json"

	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const nameAnd = "and"

type And []Criteria

func (c And) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]Criteria{
		nameAnd: c,
	})
}

func (And) criteria() {}

type definitionAnd struct{}

func (definitionAnd) name() string {
	return nameAnd
}

func (definitionAnd) compile(ctx compilerContext) (Criteria, error) {
	return compileBoolean[And](ctx, nameAnd)
}

func (definitionAnd) resultTypes() map[ResultType]struct{} {
	return map[ResultType]struct{}{
		ResultTypeContent:        {},
		ResultTypeTorrent:        {},
		ResultTypeTorrentContent: {},
		ResultTypeTorrentFile:    {},
	}
}

func (definitionAnd) JSONSchema() json_schema.JSONSchema {
	return specBoolean(nameAnd).JSONSchema()
}

func specBoolean(name string) json_spec.Typed[[]any] {
	return json_spec.SingleKeyValue[[]any]{
		Key:       name,
		ValueSpec: criteriaSpec,
	}
}

func compileBoolean[T ~[]Criteria](ctx compilerContext, name string) (T, error) {
	payload, err := json_spec.SingleKeyValue[[]any]{
		Key:       name,
		ValueSpec: criteriaSpec,
	}.
		Parse(ctx.jsonSpec)
	if err != nil {
		return nil, ctx.Error(err)
	}

	criteria, err := compileCriteriaCtx(ctx.child(name, payload))
	if err != nil {
		return nil, ctx.Fatal(err)
	}

	return T(criteria), nil
}
