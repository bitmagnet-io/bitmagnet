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

	for i := range v.NumField() {
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

	if _, ok := spec.params[helpParam]; !ok {
		spec.params[helpParam] = Param{
			Name:     helpParam,
			Abbr:     "h",
			Doc:      "Show help for this command and exit",
			Type:     paramTypeBool,
			Required: false,
		}
		spec.paramKeys = append(spec.paramKeys, helpParam)
	}

	if spec.cmd == nil {
		return spec, fmt.Errorf("%w: %w", ErrCompilation, ErrCmdNotEmbedded)
	}

	if namer, ok := cmd.(interface{ Name() string }); ok {
		spec.name = namer.Name()
	} else if spec.name == "" {
		name := strcase.ToKebab(t.Name())
		name = strings.TrimSuffix(name, "-cmd")
		name = strings.TrimSuffix(name, "-command")
		spec.name = name
	}

	if !regexParamKey.MatchString(spec.name) {
		return spec, fmt.Errorf("%w: %w: %s", ErrCompilation, ErrInvalidName, spec.name)
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
				s.name = v
			case "doc":
				s.doc = v
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
		//nolint:nilnil
		return nil, nil
	}

	paramType, err := reflectTypeToParamType(t.Type)
	if err != nil {
		return nil, err
	}

	param := Param{
		index:    t.Index,
		Name:     strcase.ToKebab(t.Name),
		Example:  strcase.ToScreamingSnake(t.Name),
		Type:     paramType,
		Multiple: paramType.isMultiple(),
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
		case "example":
			param.Example = v
		case "required":
			param.Required = true
			// case "csv":
			// 	param.CSV = true
		case "doc":
			param.Doc = v
		case "default":
			param.Default = v
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

	kvs := make(map[string]string)

	var (
		parts    []string
		current  strings.Builder
		inQuotes bool
	)

	// Split by commas, respecting single quotes

	//nolint:gocritic
	for _, ch := range tag {
		if ch == '\'' {
			inQuotes = !inQuotes

			current.WriteRune(ch)
		} else if ch == ',' && !inQuotes {
			parts = append(parts, current.String())
			current.Reset()
		} else {
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	// Parse key=value pairs
	for _, part := range parts {
		split := strings.SplitN(part, "=", 2)
		key := strings.TrimSpace(split[0])

		var value string
		if len(split) > 1 {
			value = strings.TrimSpace(split[1])
			// Remove surrounding single quotes if present
			if len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'' {
				value = value[1 : len(value)-1]
			}
		}

		kvs[key] = value
	}

	return kvs
}
