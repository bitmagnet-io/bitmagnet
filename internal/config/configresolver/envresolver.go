package configresolver

import (
	"fmt"
	"reflect"
	"strings"
)

type envResolver struct {
	baseResolver
	e map[string]string
}

func NewEnv(e map[string]string, options ...Option) Resolver {
	r := &envResolver{e: e}
	r.applyOptions(append([]Option{WithKey("env")}, options...)...)
	return r
}

func (r envResolver) Resolve(path []string, valueType reflect.Type) (interface{}, bool, error) {
	envKey := strings.ToUpper(strings.Join(path, "_"))
	envValue, ok := r.e[envKey]
	if !ok {
		return nil, false, nil
	}
	coercedValue, coerceErr := coerceStringValue(envValue, valueType)
	if coerceErr != nil {
		return nil, true, fmt.Errorf("error coercing env key '%s' with value '%s' to type %v: %w", envKey, envValue, valueType, coerceErr)
	}
	return coercedValue, true, nil
}
