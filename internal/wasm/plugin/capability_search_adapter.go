package plugin

import "github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"

const capabilitySearchAdapterName = "search_adapter"

type CapabilitySearchAdapter struct {
	Name string `json:"name"`
}

func (CapabilitySearchAdapter) capability() string {
	return capabilitySearchAdapterName
}

func (CapabilitySearchAdapter) jsonSchema() json_schema.JSONSchema {
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
