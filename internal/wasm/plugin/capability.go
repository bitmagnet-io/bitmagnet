package plugin

import "github.com/bitmagnet-io/bitmagnet/pkg/json_schema"

type Capability interface {
	capability() string
	jsonSchema() json_schema.JSONSchema
}

type Capabilities struct {
	SearchAdapter *CapabilitySearchAdapter `json:"search_adapter,omitempty"`
	Indexer       *CapabilityIndexer       `json:"indexer,omitempty"`
	HTTPHandler   *CapabilityHTTPHandler   `json:"http_handler,omitempty"`
	TorrentTarget *CapabilityTorrentTarget `json:"torrent_target,omitempty"`
}
