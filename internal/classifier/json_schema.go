package classifier

import (
	"encoding/json"
)

type JSONSchema map[string]any

func (s JSONSchema) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(map[string]any(s), "", "  ")
}

const schemaID = "https://bitmagnet.io/schemas/classifier-0.1.json"

func (f features) JSONSchema() JSONSchema {
	return map[string]any{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"$id":     schemaID,
		"type":    "object",
		"properties": map[string]any{
			"$schema": map[string]any{
				"const": schemaID,
			},
			"workflows": map[string]any{
				"type": "object",
				"additionalProperties": map[string]any{
					"$ref": "#/definitions/action",
				},
			},
			"flag_definitions": map[string]any{
				"type": "object",
				"additionalProperties": map[string]any{
					"type": "string",
					"enum": FlagTypeValues(),
				},
			},
			"flags": map[string]any{
				"type":                 "object",
				"additionalProperties": true,
			},
			"keywords": map[string]any{
				"type": "object",
				"additionalProperties": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "string",
					},
				},
			},
			"extensions": map[string]any{
				"type": "object",
				"additionalProperties": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "string",
					},
				},
			},
		},
		"additionalProperties": false,
		"definitions": func() map[string]any {
			defs := map[string]any{
				"action": map[string]any{
					"oneOf": []map[string]any{
						{
							"$ref": "#/definitions/action_single",
						},
						{
							"$ref": "#/definitions/action_multi",
						},
					},
				},
				"action_multi": map[string]any{
					"type": "array",
					"items": map[string]any{
						"$ref": "#/definitions/action_single",
					},
				},
				"action_single": map[string]any{
					"oneOf": func() []map[string]any {
						var result []map[string]any
						for _, def := range f.actions {
							result = append(result, map[string]any{
								"$ref": "#/definitions/action__" + def.name(),
							})
						}
						return result
					}(),
				},
				"condition": map[string]any{
					"oneOf": func() []map[string]any {
						var result []map[string]any
						for _, def := range f.conditions {
							result = append(result, map[string]any{
								"$ref": "#/definitions/condition__" + def.name(),
							})
						}
						return result
					}(),
				},
			}
			for _, def := range f.actions {
				defs["action__"+def.name()] = def.JSONSchema()
			}
			for _, def := range f.conditions {
				defs["condition__"+def.name()] = def.JSONSchema()
			}
			return defs
		}(),
	}
}

func DefaultJSONSchema() JSONSchema {
	return defaultFeatures.JSONSchema()
}
