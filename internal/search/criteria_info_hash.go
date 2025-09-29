package search

import (
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

const nameInfoHash = "infoHash"

type CriteriaInfoHash []protocol.ID

func (CriteriaInfoHash) criteria() {}

var criteriaSpecInfoHash = json_spec.SingleKeyValue[[]protocol.ID]{
	Key: nameInfoHash,
	ValueSpec: json_spec.MustSucceed[[]protocol.ID]{
		Typed: json_spec.List[protocol.ID]{
			ItemSpec: json_spec.Generic[protocol.ID]{
				Schema: json_schema.MustNew(
					json_schema.Typed(json_schema.TypeString),
				),
				Parser: func(ctx json_spec.ParseContext) (protocol.ID, error) {
					str, ok := ctx.Source.(string)
					if !ok {
						return protocol.ID{}, fmt.Errorf("expected string, got %T", ctx.Source)
					}

					return protocol.ParseID(str)
				},
			},
		},
	},
}

type definitionInfoHash struct{}

func (definitionInfoHash) name() string {
	return nameInfoHash
}

func (definitionInfoHash) compile(ctx compilerContext) (Criteria, error) {
	infoHashes, err := criteriaSpecInfoHash.Parse(ctx.jsonSpec)
	if err != nil {
		return nil, ctx.Error(err)
	}

	return CriteriaInfoHash(infoHashes), nil
}

func (definitionInfoHash) resultTypes() map[ResultType]struct{} {
	return map[ResultType]struct{}{
		ResultTypeTorrent:        {},
		ResultTypeTorrentContent: {},
		ResultTypeTorrentFile:    {},
	}
}

func (definitionInfoHash) JSONSchema() json_schema.JSONSchema {
	return criteriaSpecInfoHash.JSONSchema()
}
