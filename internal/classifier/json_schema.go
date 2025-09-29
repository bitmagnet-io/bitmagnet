package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

const schemaID = "https://bitmagnet.io/schemas/classifier-0.1.json"

func (f features) JSONSchema() json_schema.JSONSchema {
	return json_schema.MustNew(
		json_schema.MetaschemaDraft7(),
		json_schema.ID(schemaID),
		json_schema.Typed(json_schema.TypeObject),
		json_schema.Properties(map[string]json_schema.JSONSchema{
			"$schema": json_schema.MustNew(
				json_schema.Const(json_schema.MustNewValue(schemaID)),
			),
			"workflows": json_schema.MustNew(
				json_schema.Typed(json_schema.TypeObject),
				json_schema.AdditionalPropertiesType(
					json_schema.MustNew(json_schema.RefDefinition("action")),
				),
			),
			"flag_definitions": json_schema.MustNew(
				json_schema.Typed(json_schema.TypeObject),
				json_schema.AdditionalPropertiesType(
					json_schema.MustNew(
						json_schema.Typed(json_schema.TypeString),
						json_schema.Enum(slice.Map(FlagTypeValues(), func(ft FlagType) json_schema.JSONValue {
							return json_schema.MustNewValue(string(ft))
						})...),
					),
				),
			),
			"flags": json_schema.MustNew(
				json_schema.Typed(json_schema.TypeObject),
				json_schema.AdditionalPropertiesTrue(),
			),
			"keywords": json_schema.MustNew(
				json_schema.Typed(json_schema.TypeObject),
				json_schema.AdditionalPropertiesType(
					json_schema.MustNew(
						json_schema.Typed(json_schema.TypeArray),
						json_schema.Items(
							json_schema.MustNew(json_schema.Typed(json_schema.TypeString)),
						),
					),
				),
			),
			"extensions": json_schema.MustNew(
				json_schema.Typed(json_schema.TypeObject),
				json_schema.AdditionalPropertiesType(
					json_schema.MustNew(
						json_schema.Typed(json_schema.TypeArray),
						json_schema.Items(
							json_schema.MustNew(json_schema.Typed(json_schema.TypeString)),
						),
					),
				),
			),
		}),
		json_schema.AdditionalPropertiesFalse(),
		json_schema.Definitions(
			map[string]json_schema.JSONSchema{
				"action": json_schema.MustNew(
					json_schema.OneOf(
						json_schema.MustNew(
							json_schema.RefDefinition("action_single"),
						),
						json_schema.MustNew(
							json_schema.RefDefinition("action_multi"),
						),
					),
				),
				"action_multi": json_schema.MustNew(
					json_schema.Typed(json_schema.TypeArray),
					json_schema.Items(
						json_schema.MustNew(json_schema.RefDefinition("action_single")),
					),
				),
				"action_single": json_schema.MustNew(
					json_schema.OneOf(slice.Map(f.actions, func(a actionDefinition) json_schema.JSONSchema {
						return json_schema.MustNew(json_schema.RefDefinition("action__" + a.name()))
					})...),
				),
				"condition": json_schema.MustNew(
					json_schema.OneOf(slice.Map(f.conditions, func(c conditionDefinition) json_schema.JSONSchema {
						return json_schema.MustNew(json_schema.RefDefinition("condition__" + c.name()))
					})...),
				),
			},
		),
		json_schema.Options(slice.Map(f.actions, func(a actionDefinition) json_schema.Option {
			return json_schema.Definition("action__"+a.name(), a.JSONSchema())
		})...),
		json_schema.Options(slice.Map(f.conditions, func(c conditionDefinition) json_schema.Option {
			return json_schema.Definition("condition__"+c.name(), c.JSONSchema())
		})...),
	)
}

func DefaultJSONSchema() json_schema.JSONSchema {
	return defaultFeatures.JSONSchema()
}
