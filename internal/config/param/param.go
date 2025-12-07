package param

import (
	"encoding/csv"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"gopkg.in/yaml.v3"
)

type Untyped interface {
	Description() string
	ValidateAny(any) error
	ParseAny(string) (any, error)
	StringifyAny(any) (string, error)
	EncodeYAMLAny(any) (yaml.Node, error)
	EncodeYAMLAnyAny(any) (any, error)
	DecodeYAMLAny(yaml.Node) (any, error)
	DecodeYAMLAnyAny(any) (any, error)
	NewDefaultAny() any
	HasExplicitDefault() bool
	ReflectType() reflect.Type
	DynamicType() (reflect.Type, bool)
	IsDynamic() bool
	JSONSchema() json_schema.JSONSchema
}

type Param[T any] interface {
	Untyped
	Validators() Validators[T]
	Validate(T) error
	EncodeYAML(T) (yaml.Node, error)
	DecodeYAML(yaml.Node) (T, error)
	Stringify(T) string
	Parse(string) (T, error)
	Equals(T, T) bool
	EnumValues() []T
	NewDefault() T
}

type param[T any] struct {
	description        string
	newDefault         func() T
	hasExplicitDefault bool
	yamlEncoder        func(T) (yaml.Node, error)
	yamlDecoder        func(yaml.Node) (T, error)
	validators         Validators[T]
	comparator         func(T, T) bool
	stringifier        func(T) string
	parser             func(string) (T, error)
	enumValues         []T
	dynamicType        reflect.Type
	jsonSchema         json_schema.JSONSchema
}

func New[T any](opts ...Option[T]) (Param[T], error) {
	p := param[T]{
		newDefault:  newDefaultZero[T],
		yamlEncoder: yamlEncoder[T],
		yamlDecoder: yamlDecoder[T],
		stringifier: stringifierSimple[T],
		comparator:  comparatorReflect[T],
		jsonSchema:  json_schema.MustNew(json_schema.Typed(json_schema.TypeString)),
	}

	for _, opt := range opts {
		if err := opt(&p); err != nil {
			return p, err
		}
	}

	if p.parser == nil {
		p.parser = parserYAML(p.yamlDecoder)
	}

	if p.hasExplicitDefault && p.jsonSchema.Default == nil {
		yamlDefault, err := p.EncodeYAMLAny(p.NewDefault())
		if err != nil {
			return p, err
		}

		err = JSONSchemaOption[T](json_schema.Default(json_schema.JSONValue(yamlDefault)))(&p)
		if err != nil {
			return p, err
		}
	}

	return p, nil
}

func MustNew[T any](opts ...Option[T]) Param[T] {
	value, err := New(opts...)
	if err != nil {
		panic(err)
	}

	return value
}

func (p param[T]) Equals(a, b T) bool {
	if p.comparator == nil {
		return false
	}

	return p.comparator(a, b)
}

func (p param[T]) Stringify(val T) string {
	return p.stringifier(val)
}

func (p param[T]) StringifyAny(val any) (string, error) {
	typed, ok := val.(T)
	if !ok {
		return "", errors.New("invalid type")
	}

	return p.Stringify(typed), nil
}

func (p param[T]) Parse(s string) (T, error) {
	return p.parser(s)
}

func (p param[T]) ParseAny(s string) (any, error) {
	return p.Parse(s)
}

func (p param[T]) Validators() Validators[T] {
	return p.validators
}

func (p param[T]) Validate(val T) error {
	return p.validators.Validate(val)
}

func (p param[T]) ValidateAny(val any) error {
	typed, ok := val.(T)
	if !ok {
		return fmt.Errorf("failed to cast %T to %T", val, typed)
	}

	return p.Validate(typed)
}

func (p param[T]) EncodeYAML(val T) (yaml.Node, error) {
	return p.yamlEncoder(val)
}

func (p param[T]) EncodeYAMLAny(val any) (yaml.Node, error) {
	typed, ok := val.(T)
	if !ok {
		return yaml.Node{}, fmt.Errorf("failed to cast %T to %T", val, typed)
	}

	return p.EncodeYAML(typed)
}

func (p param[T]) EncodeYAMLAnyAny(val any) (any, error) {
	node, err := p.EncodeYAMLAny(val)
	if err != nil {
		return nil, err
	}

	var decoded any
	err = node.Decode(&decoded)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func (p param[T]) DecodeYAML(node yaml.Node) (T, error) {
	return p.yamlDecoder(node)
}

func (p param[T]) DecodeYAMLAny(node yaml.Node) (any, error) {
	return p.DecodeYAML(node)
}

func (p param[T]) DecodeYAMLAnyAny(value any) (any, error) {
	var node yaml.Node

	if jsonSchemaValue, ok := value.(json_schema.JSONValue); ok {
		node = yaml.Node(jsonSchemaValue)
	} else if valueNode, ok := value.(yaml.Node); ok {
		node = valueNode
	} else {
		err := node.Encode(value)
		if err != nil {
			return nil, err
		}
	}

	return p.DecodeYAMLAny(node)
}

func (p param[T]) ReflectType() reflect.Type {
	var t T
	return reflect.TypeOf(t)
}

func (p param[T]) DynamicType() (reflect.Type, bool) {
	return p.dynamicType, p.dynamicType != nil
}

func (p param[T]) IsDynamic() bool {
	return p.dynamicType != nil
}

func (p param[T]) EnumValues() []T {
	return p.enumValues
}

func (p param[T]) JSONSchema() json_schema.JSONSchema {
	return p.jsonSchema
}

func (p param[T]) Description() string {
	return p.description
}

func (p param[T]) NewDefault() T {
	return p.newDefault()
}

func (p param[T]) NewDefaultAny() any {
	return p.NewDefault()
}

func (p param[T]) HasExplicitDefault() bool {
	return p.hasExplicitDefault
}

func newDefaultZero[T any]() T {
	var zero T
	return zero
}

func parserYAML[T any](elementParser func(yaml.Node) (T, error)) func(string) (T, error) {
	return func(s string) (T, error) {
		var node yaml.Node
		if err := yaml.Unmarshal([]byte(s), &node); err != nil {
			return *new(T), err
		}
		return elementParser(node)
	}
}

func parserSlice[E any, T ~[]E](elementParser func(string) (E, error)) func(string) (T, error) {
	return func(s string) (T, error) {
		parts, err := csv.NewReader(strings.NewReader(s)).Read()
		if err != nil {
			return nil, err
		}
		var result T
		for _, item := range parts {
			parsed, err := elementParser(item)
			if err != nil {
				return nil, err
			}
			result = append(result, parsed)
		}
		return result, nil
	}
}

func parserDynamic[T any](elementParser func(string) (T, error)) func(string) (*atomic.Value[T], error) {
	return func(s string) (*atomic.Value[T], error) {
		element, err := elementParser(s)
		if err != nil {
			return nil, err
		}

		return atomic.NewValue(element), nil
	}
}

func stringifierSlice[E any, T ~[]E](elementStringifier func(E) string) func(T) string {
	return func(sl T) string {
		var b strings.Builder
		csv.NewWriter(&b).Write(slice.Map(sl, elementStringifier))
		return b.String()
	}
}

func stringifierDynamic[T any](elementStringifier func(T) string) func(*atomic.Value[T]) string {
	return func(value *atomic.Value[T]) string {
		return elementStringifier(value.Get())
	}
}
