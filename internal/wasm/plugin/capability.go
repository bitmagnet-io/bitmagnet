package plugin

import "github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"

type Capability interface {
	capability() string
	jsonSchema() json_schema.JSONSchema
}

type Capabilities struct {
	SearchAdapter *CapabilitySearchAdapter `json:"search_adapter,omitempty"`
	Indexer       *CapabilityIndexer       `json:"indexer,omitempty"`
	HTTPHandler   *CapabilityHTTPHandler   `json:"http_handler,omitempty"`
	Receiver      *CapabilityReceiver      `json:"receiver,omitempty"`
}
