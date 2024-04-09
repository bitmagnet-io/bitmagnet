package classifier

import "encoding/json"

type JsonSchema map[string]any

func (s JsonSchema) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(map[string]any(s), "", "  ")
}

func (f features) JsonSchema() JsonSchema {
	return map[string]any{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"$id":     "https://bitmagnet.io/schemas/classifier-0.1.json",
		"type":    "object",
		"properties": map[string]any{
			"workflows": map[string]any{
				"type": "object",
				"additionalProperties": map[string]any{
					"$ref": "#/$defs/action",
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
		"$defs": func() map[string]any {
			defs := map[string]any{
				"action": map[string]any{
					"oneOf": []map[string]any{
						{
							"$ref": "#/$defs/action_single",
						},
						{
							"$ref": "#/$defs/action_multi",
						},
					},
				},
				"action_multi": map[string]any{
					"type": "array",
					"items": map[string]any{
						"$ref": "#/$defs/action_single",
					},
				},
				"action_single": map[string]any{
					"oneOf": func() []map[string]any {
						var result []map[string]any
						for _, def := range f.actions {
							result = append(result, map[string]any{
								"$ref": "#/$defs/action__" + def.name(),
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
								"$ref": "#/$defs/condition__" + def.name(),
							})
						}
						return result
					}(),
				},
			}
			for _, def := range f.actions {
				defs["action__"+def.name()] = def.JsonSchema()
			}
			for _, def := range f.conditions {
				defs["condition__"+def.name()] = def.JsonSchema()
			}
			return defs
		}(),
	}
}

func DefaultJsonSchema() JsonSchema {
	return defaultFeatures.JsonSchema()
}
