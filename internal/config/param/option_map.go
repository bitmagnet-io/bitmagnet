package param

// func Map[T any](params map[string]Untyped) Option[T] {
// 	return Options(
// 		YAMLDecoder(func(node yaml.Node) (T, error) {
// 			var raw map[string]yaml.Node
// 			if err := node.Decode(&raw); err != nil {
// 				var zero T
// 				return zero, err
// 			}

// 			result := make(map[string]any)
// 			for key, param := range params {
// 				if paramNode, ok := raw[key]; ok {
// 					value, err := param.DecodeYAMLAny(paramNode)
// 					if err != nil {
// 						var zero T
// 						return zero, err
// 					}
// 					result[key] = value
// 				}
// 			}

// 			var value T
// 			decoder, err := newDecoder(&value)
// 			if err != nil {
// 				var zero T
// 				return zero, err
// 			}
// 			if err := decoder.Decode(result); err != nil {
// 				var zero T
// 				return zero, err
// 			}
// 			return value, nil
// 		}),
// 		JSONSchemaOption[T](
// 			json_schema.Typed(json_schema.TypeObject),
// 			json_schema.Properties(
// 				func() map[string]json_schema.JSONSchema {
// 					props := make(map[string]json_schema.JSONSchema)
// 					for key, param := range params {
// 						props[key] = param.JSONSchema()
// 					}
// 					return props
// 				}(),
// 			),
// 			json_schema.Required(json_schema.RequiredFields(
// 				slice.Filter(slices.Sorted(maps.Keys(params)), func(key string) bool {
// 					req := params[key].JSONSchema().Required
//           typedReq, ok := any(req).(*json_schema.RequiredFields)
//           return ok && typedReq != nil && *typedReq {

//           }
// 				}),
// 			)),
// 		),
// 	)
// }
