package search

import (
	"encoding/json"

	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const nameOr = "or"

type Or []Criteria

func (c Or) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]Criteria{
		nameOr: c,
	})
}

func (Or) criteria() {}

type definitionOr struct{}

func (definitionOr) name() string {
	return nameOr
}

func (definitionOr) compile(ctx compilerContext) (Criteria, error) {
	return compileBoolean[Or](ctx, nameOr)
}

func (definitionOr) resultTypes() map[ResultType]struct{} {
	return map[ResultType]struct{}{
		ResultTypeContent:        {},
		ResultTypeTorrent:        {},
		ResultTypeTorrentContent: {},
		ResultTypeTorrentFile:    {},
	}
}

func (definitionOr) JSONSchema() json_schema.JSONSchema {
	return specBoolean(nameOr).JSONSchema()
}
