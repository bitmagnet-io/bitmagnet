package classifier

import (
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protobuf"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

type flagDefinitions map[string]FlagType

func (d flagDefinitions) merge(other flagDefinitions) (flagDefinitions, error) {
	result := make(flagDefinitions)

	for k, v := range d {
		tp, ok := other[k]
		if ok && tp != v {
			return nil, fmt.Errorf("conflicting flag definition %s", k)
		}

		result[k] = v
	}

	for k, v := range other {
		if _, ok := result[k]; !ok {
			result[k] = v
		}
	}

	return result, nil
}

type Flags map[string]any

func (f Flags) merge(other Flags) Flags {
	result := make(Flags)

	for k, v := range f {
		if _, ok := other[k]; ok {
			result[k] = other[k]
		} else {
			result[k] = v
		}
	}

	for k, v := range other {
		if _, ok := result[k]; !ok {
			result[k] = v
		}
	}

	return result
}

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
				strVal, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("could not convert type %T to string", v)
				}

				nativeVal[i] = strVal
			}

			return types.NewStringList(types.DefaultTypeAdapter, nativeVal), nil
		}
	case FlagTypeContentTypeList:
		if sliceVal, ok := rawVal.([]any); ok {
			celVal := make([]protobuf.Classification_ContentType, len(sliceVal))

			for i, v := range sliceVal {
				strVal, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("could not convert type %T to content type", v)
				}

				var ct model.NullContentType

				if strVal != "unknown" {
					parsed, parseErr := model.ParseContentType(strVal)
					if parseErr != nil {
						return nil, fmt.Errorf(
							"could not parse content type %s: %w",
							strVal,
							parseErr,
						)
					}

					ct = model.NewNullContentType(parsed)
				}

				celVal[i] = protobuf.NewContentType(ct)
			}

			return types.NewDynamicList(types.DefaultTypeAdapter, celVal), nil
		}
	default:
		return nil, ErrInvalidFlagType
	}

	return nil, fmt.Errorf("could not convert type %T to %s", rawVal, t)
}

type compiledFlags map[string]ref.Val
