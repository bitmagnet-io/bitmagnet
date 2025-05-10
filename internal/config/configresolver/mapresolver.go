package configresolver

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/go-viper/mapstructure/v2"
	"github.com/iancoleman/strcase"
)

type mapResolver struct {
	baseResolver
	validator *validator.Validate
	m         map[string]interface{}
}

func NewMap(m map[string]interface{}, val *validator.Validate, options ...Option) Resolver {
	r := &mapResolver{m: m, validator: val}
	r.applyOptions(append([]Option{WithKey("map")}, options...)...)

	return r
}

func (r mapResolver) Resolve(path []string, valueType reflect.Type) (any, bool, error) {
	if len(path) == 0 {
		return r.m, true, nil
	}

	v := r.m

	currentPath := make([]string, 0, len(path))
	for i, p := range path {
		currentPath = append(currentPath, p)

		rawV, rawOk := v[p]
		if !rawOk {
			break
		}

		if i >= len(path)-1 {
			if strV, strOk := rawV.(string); strOk {
				coerced, coerceErr := coerceStringValue(strV, valueType)
				if coerceErr != nil {
					return nil, true, fmt.Errorf(
						"error coercing configcmd map path '%s' with value '%s' to type %v: %w",
						currentPath, strV, valueType, coerceErr)
				}

				return coerced, true, nil
			} else if sliceV, sliceOk := rawV.([]interface{}); sliceOk {
				resolvedSlice, err := r.resolveSlice(currentPath, sliceV, valueType)
				return resolvedSlice, true, err
			}

			return rawV, true, nil
		}

		mapV, mapOk := rawV.(map[string]interface{})
		if !mapOk {
			return nil, true, fmt.Errorf(
				"expected map[string]interface{} at path %v, got %T",
				currentPath,
				rawV,
			)
		}

		v = mapV
	}

	return nil, false, nil
}

func (r mapResolver) resolveSlice(currentPath []string, sliceV []any, valueType reflect.Type) (any, error) {
	if valueType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("received slice at path '%s', expected %s", currentPath, valueType.String())
	}

	resolvedSlice := make([]any, 0, len(sliceV))

	for _, sliceItem := range sliceV {
		resolvedValue := reflect.New(valueType.Elem())

		decoder, decoderErr := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result: resolvedValue.Interface(),
			MatchName: func(mapKey, fieldName string) bool {
				return mapKey == strcase.ToSnake(fieldName)
			},
		})
		if decoderErr != nil {
			return nil, decoderErr
		}

		if decodeErr := decoder.Decode(sliceItem); decodeErr != nil {
			return nil, decodeErr
		}

		if valueType.Elem().Kind() == reflect.Struct {
			if validateErr := r.validator.Struct(resolvedValue.Interface()); validateErr != nil {
				return nil, validateErr
			}
		}

		resolvedSlice = append(resolvedSlice, reflect.Indirect(resolvedValue).Interface())
	}

	return resolvedSlice, nil
}
