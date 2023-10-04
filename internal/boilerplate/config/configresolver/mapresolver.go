package configresolver

import (
	"fmt"
	"reflect"
)

type mapResolver struct {
	baseResolver
	m map[string]interface{}
}

func NewMap(m map[string]interface{}, options ...Option) Resolver {
	r := &mapResolver{m: m}
	r.applyOptions(append([]Option{WithKey("map")}, options...)...)
	return r
}

func (r mapResolver) Resolve(path []string, valueType reflect.Type) (interface{}, bool, error) {
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
		if i < len(path)-1 {
			mapV, mapOk := rawV.(map[string]interface{})
			if !mapOk {
				return nil, true, fmt.Errorf("expected map[string]interface{} at path %v, got %T", currentPath, rawV)
			}
			v = mapV
		} else {
			if strV, strOk := rawV.(string); strOk {
				coerced, coerceErr := coerceStringValue(strV, valueType)
				if coerceErr != nil {
					return nil, true, fmt.Errorf("error coercing config map path '%s' with value '%s' to type %v: %w", currentPath, strV, valueType, coerceErr)
				}
				return coerced, true, nil
			}
			return rawV, true, nil
		}
	}
	return nil, false, nil
}
