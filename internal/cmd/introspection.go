package cmd

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

func introspect(cmd Command) (Spec, error) {
	v := reflect.ValueOf(cmd)
	if v.IsZero() {
		return Spec{}, fmt.Errorf("%w: %w", ErrCompilation, ErrZeroValue)
	}

	t := v.Type()
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	if t.Kind() != reflect.Struct {
		return Spec{}, fmt.Errorf("%w: %w", ErrCompilation, ErrNonStructValue)
	}

	spec := Spec{
		params: make(map[string]Param),
	}

	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)

		if !ft.IsExported() {
			continue
		}

		switch {
		case ft.Type == reflect.TypeOf(Cmd{}):
			err := applyCmd(ft)(&spec)
			if err != nil {
				return spec, fmt.Errorf("%w: %w", ErrCompilation, err)
			}
		case ft.Anonymous:
		default:
			err := applyParam(ft)(&spec)
			if err != nil {
				return spec, fmt.Errorf("%w: %w", ErrCompilation, err)
			}
		}
	}

	if spec.cmd == nil {
		return spec, fmt.Errorf("%w: %w", ErrCompilation, ErrCmdNotEmbedded)
	}

	if namer, ok := cmd.(interface{ Name() string }); ok {
		spec.Name = namer.Name()
	} else {
		name := strcase.ToKebab(t.Name())
		name = strings.TrimSuffix(name, "-cmd")
		name = strings.TrimSuffix(name, "-command")
		spec.Name = name
	}

	if !regexParamKey.MatchString(spec.Name) {
		return spec, fmt.Errorf("%w: %w: %s", ErrCompilation, ErrInvalidName, spec.Name)
	}

	return spec, nil
}

func applyCmd(t reflect.StructField) func(*Spec) error {
	return func(s *Spec) error {
		if !t.Anonymous {
			return ErrCmdNotEmbedded
		}

		kvs := tagKeyValues(t.Tag.Get("cmd"))

		for k, v := range kvs {
			switch k {
			case "name":
				s.Name = v
			default:
				return fmt.Errorf("%w: %s", ErrUnknownCmdTag, k)
			}
		}

		s.cmd = t.Index

		return nil
	}
}

func applyParam(t reflect.StructField) func(*Spec) error {
	return func(s *Spec) error {
		param, err := extractParam(t)
		if err != nil {
			return err
		}
		if param != nil {
			s.params[param.Name] = *param
			s.paramKeys = append(s.paramKeys, param.Name)
		}
		return nil
	}
}

func extractParam(t reflect.StructField) (*Param, error) {
	tag := t.Tag.Get("cmd")
	if tag == "-" {
		return nil, nil
	}

	paramType, err := reflectTypeToParamType(t.Type)
	if err != nil {
		return nil, err
	}

	param := Param{
		index:       t.Index,
		Name:        strcase.ToKebab(t.Name),
		Placeholder: strcase.ToScreamingSnake(t.Name),
		Type:        paramType,
		Multiple:    paramType.isMultiple(),
	}

	tagValues := tagKeyValues(tag)

	for k, v := range tagValues {
		switch k {
		case "param":
			if !regexParamKey.MatchString(v) {
				return nil, fmt.Errorf("%w: %s: %s", ErrInvalidName, t.Name, v)
			}
			param.Name = v
		case "abbr":
			if !regexParamKey.MatchString(v) {
				return nil, fmt.Errorf("%w: %s: %s", ErrInvalidAbbr, t.Name, v)
			}
			param.Abbr = v
		case "placeholder":
			param.Placeholder = v
		case "required":
			param.Required = true
		// case "csv":
		// 	param.CSV = true
		default:
			return nil, fmt.Errorf("%w: %s: %s", ErrUnknownParamTag, t.Name, k)
		}
	}

	return &param, nil
}

var textUnmarshalerType = reflect.ValueOf(new(encoding.TextUnmarshaler)).Type().Elem()

func reflectTypeToParamType(t reflect.Type) (paramType, error) {
	if reflect.PointerTo(t).Implements(textUnmarshalerType) {
		return paramTypeTextUnmarshaler, nil
	}
	switch t.Kind() {
	case reflect.String:
		return paramTypeString, nil
	case reflect.Bool:
		return paramTypeBool, nil
	case reflect.Int:
		return paramTypeInt, nil
	default:
		return paramTypeUnknown, fmt.Errorf("%w: %s", ErrUnsupportedParamType, t.String())
	}
}

func tagKeyValues(tag string) map[string]string {
	if tag == "" {
		return nil
	}
	parts := strings.Split(tag, ",")
	kvs := make(map[string]string, len(parts))
	for _, part := range parts {
		split := strings.SplitN(part, "=", 2)
		key := split[0]
		var value string
		if len(split) > 1 {
			value = split[1]
		}
		kvs[key] = value
	}
	return kvs
}
