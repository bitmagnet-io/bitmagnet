package workflow

import (
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protobuf"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

type flagDefinitionPayload struct {
	Type    FlagType
	Default interface{}
}

type flagDefinitions map[string]flagDefinitionPayload

type flags map[string]any

func (t FlagType) celType() *cel.Type {
	switch t {
	case FlagTypeBool:
		return cel.BoolType
	case FlagTypeString:
		return cel.StringType
	case FlagTypeInt:
		return cel.IntType
	case FlagTypeStringList:
		return cel.ListType(cel.StringType)
	case FlagTypeContentTypeList:
		return cel.ListType(cel.IntType)
	default:
		return nil
	}
}

func (t FlagType) celVal(rawVal any) (ref.Val, error) {
	switch t {
	case FlagTypeBool:
		if nativeVal, ok := rawVal.(bool); ok {
			return types.Bool(nativeVal), nil
		}
	case FlagTypeString:
		if nativeVal, ok := rawVal.(string); ok {
			return types.String(nativeVal), nil
		}
	case FlagTypeInt:
		if nativeVal, ok := rawVal.(int); ok {
			return types.Int(nativeVal), nil
		}
	case FlagTypeStringList:
		if sliceVal, ok := rawVal.([]any); ok {
			nativeVal := make([]string, len(sliceVal))
			for i, v := range sliceVal {
				if strVal, ok := v.(string); ok {
					nativeVal[i] = strVal
				} else {
					return nil, fmt.Errorf("could not convert type %T to string", v)
				}
			}
			return types.NewStringList(types.DefaultTypeAdapter, nativeVal), nil
		}
	case FlagTypeContentTypeList:
		if sliceVal, ok := rawVal.([]any); ok {
			celVal := make([]protobuf.Classification_ContentType, len(sliceVal))
			for i, v := range sliceVal {
				if strVal, ok := v.(string); ok {
					var ct model.NullContentType
					if strVal != "unknown" {
						if parsed, parseErr := model.ParseContentType(strVal); parseErr != nil {
							return nil, fmt.Errorf("could not parse content type %s: %w", strVal, parseErr)
						} else {
							ct = model.NewNullContentType(parsed)
						}
					}
					celVal[i] = protobuf.NewContentType(ct)
				} else {
					return nil, fmt.Errorf("could not convert type %T to content type", v)
				}
			}
			return types.NewDynamicList(types.DefaultTypeAdapter, celVal), nil
		}
	default:
		return nil, ErrInvalidFlagType
	}
	return nil, fmt.Errorf("could not convert type %T to %s", rawVal, t)
}
