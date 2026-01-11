package plugin

import "github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"

type Permission interface {
	permission()
	jsonSchema() json_schema.JSONSchema
}

type Permissions struct {
	FS   *PermissionFS
	HTTP *PermissionHTTP
}
