package plugin

import "github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"

const capabilityReceiverName = "receiver"

type CapabilityReceiver struct {
	Name string `json:"name"`
}

func (CapabilityReceiver) capability() string {
	return capabilityReceiverName
}

func (CapabilityReceiver) jsonSchema() json_schema.JSONSchema {
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
