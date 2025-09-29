package search

import "github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"

const nameNot = "not"

type Not []Criteria

func (Not) criteria() {}

type definitionNot struct{}

func (definitionNot) name() string {
	return nameNot
}

func (definitionNot) compile(ctx compilerContext) (Criteria, error) {
	return compileBoolean[Not](ctx, nameNot)
}

func (definitionNot) resultTypes() map[ResultType]struct{} {
	return map[ResultType]struct{}{
		ResultTypeContent:        {},
		ResultTypeTorrent:        {},
		ResultTypeTorrentContent: {},
		ResultTypeTorrentFile:    {},
	}
}

func (definitionNot) JSONSchema() json_schema.JSONSchema {
	return specBoolean(nameNot).JSONSchema()
}
