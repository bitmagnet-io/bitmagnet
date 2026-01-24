package plugin

import "github.com/bitmagnet-io/bitmagnet/pkg/json_schema"

const capabilityIndexerName = "indexer"

type CapabilityIndexer struct {
	Name string `json:"name"`
}

func (CapabilityIndexer) capability() string {
	return capabilityIndexerName
}

func (CapabilityIndexer) jsonSchema() json_schema.JSONSchema {
	return json_schema.MustNew(
		json_schema.Typed(json_schema.TypeObject),
		json_schema.Properties(map[string]json_schema.JSONSchema{
			"name": json_schema.MustNew(
				json_schema.Typed(json_schema.TypeString),
			),
		}),
		json_schema.Required(json_schema.RequiredFields{
			"name",
		}),
		json_schema.AdditionalPropertiesFalse(),
	)
}
