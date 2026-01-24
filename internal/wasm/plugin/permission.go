package plugin

import "github.com/bitmagnet-io/bitmagnet/pkg/json_schema"

type Permission interface {
	permission()
	jsonSchema() json_schema.JSONSchema
}

type Permissions struct {
	FS   *PermissionFS   `json:"fs,omitempty"`
	HTTP *PermissionHTTP `json:"http,omitempty"`
}
