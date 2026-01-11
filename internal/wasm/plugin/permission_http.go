package plugin

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/wasm/host/http_client"
)

type PermissionHTTP struct{}

func (PermissionHTTP) permission() {}

func (PermissionHTTP) jsonSchema() json_schema.JSONSchema {
	return json_schema.MustNew(
		json_schema.Typed(json_schema.TypeObject),
	)
}

func (p PermissionHTTP) build(b *instanceBuilder) {
	b.instantiators = append(b.instantiators, http_client.Instantiator())
}
