package cmd

import (
	"fmt"
	"reflect"

	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type Spec struct {
	Name      string
	cmd       []int
	params    map[string]Param
	paramKeys []string
}

func (s Spec) Params() []Param {
	return slice.Map(s.paramKeys, func(key string) Param {
		return s.params[key]
	})
}

func (s Spec) param(key string) (Param, error) {
	param, ok := s.params[key]
	if !ok {
		return param, fmt.Errorf("%w: %s", ErrUnknownParam, key)
	}

	return param, nil
}

func (s Spec) newInstance(cmd Command, args []string) *instance {
	reflectValue := reflect.ValueOf(cmd).Elem()
	reflectValues := make(map[string]reflect.Value)
	for name, param := range s.params {
		reflectValues[name] = reflectValue.FieldByIndex(param.index)
	}

	return &instance{
		Spec:          s,
		Command:       cmd,
		cmdValue:      reflectValue.FieldByIndex(s.cmd),
		values:        make(map[string][]any),
		reflectValues: reflectValues,
		args:          args,
	}
}

type paramType int

const (
	paramTypeUnknown paramType = iota
	paramTypeString
	paramTypeBool
	paramTypeInt
	// paramTypeStringSlice
	paramTypeTextUnmarshaler
)

func (p paramType) isMultiple() bool {
	// switch p {
	// case paramTypeStringSlice:
	// 	return true
	// default:
	return false
	// }
}

type Param struct {
	Name        string
	Abbr        string
	Placeholder string
	Type        paramType
	Required    bool
	Multiple    bool
	// CSV         bool

	index []int
	// validate    string
	// trailer   bool
}
