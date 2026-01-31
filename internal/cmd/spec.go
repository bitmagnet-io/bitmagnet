package cmd

import (
	"reflect"

	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type Spec struct {
	name      string
	doc       string
	cmd       []int
	params    map[string]Param
	paramKeys []string
}

func (s Spec) Params() []Param {
	return slice.Map(s.paramKeys, func(key string) Param {
		return s.params[key]
	})
}

func (s Spec) param(key string) (Param, bool) {
	param, ok := s.params[key]

	return param, ok
}

func (s Spec) newInstance(cmd Command, args []string) *instance {
	reflectValue := reflect.ValueOf(cmd).Elem()
	reflectValues := make(map[string]reflect.Value)

	for name, param := range s.params {
		if param.index != nil {
			reflectValues[name] = reflectValue.FieldByIndex(param.index)
		}
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

func (paramType) isMultiple() bool {
	// switch p {
	// case paramTypeStringSlice:
	// 	return true
	// default:
	return false
	// }
}

type Param struct {
	Name     string
	Abbr     string
	Example  string
	Default  string
	Doc      string
	Type     paramType
	Required bool
	Multiple bool

	index []int
}
