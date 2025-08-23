package param

import (
	"encoding/csv"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"gopkg.in/yaml.v3"
)

type Untyped interface {
	Doc() string
	ValidateAny(any) error
	ParseAny(string) (any, error)
	StringifyAny(any) (string, error)
	EncodeYAMLAny(any) (yaml.Node, error)
	DecodeYAMLAny(yaml.Node) (any, error)
	NewDefaultAny() any
	ReflectType() reflect.Type
	DynamicType() (reflect.Type, bool)
	IsDynamic() bool
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
	doc         string
	newDefault  func() T
	yamlEncoder func(T) (yaml.Node, error)
	yamlDecoder func(yaml.Node) (T, error)
	validators  Validators[T]
	comparator  func(T, T) bool
	stringifier func(T) string
	parser      func(string) (T, error)
	enumValues  []T
	dynamicType reflect.Type
}

func New[T any](opts ...Option[T]) (Param[T], error) {
	v := param[T]{
		newDefault:  newDefaultZero[T],
		yamlEncoder: yamlEncoder[T],
		yamlDecoder: yamlDecoder[T],
		stringifier: stringifierSimple[T],
		comparator:  comparatorReflect[T],
	}
	for _, opt := range opts {
		opt(&v)
	}

	if v.parser == nil {
		v.parser = parserYAML(v.yamlDecoder)
	}

	return v, nil
}

func MustNew[T any](opts ...Option[T]) Param[T] {
	value, err := New(opts...)
	if err != nil {
		panic(err)
	}

	return value
}

func (v param[T]) Equals(a, b T) bool {
	if v.comparator == nil {
		return false
	}

	return v.comparator(a, b)
}

func (v param[T]) Stringify(val T) string {
	return v.stringifier(val)
}

func (v param[T]) StringifyAny(val any) (string, error) {
	typed, ok := val.(T)
	if !ok {
		return "", errors.New("invalid type")
	}

	return v.Stringify(typed), nil
}

func (v param[T]) Parse(s string) (T, error) {
	return v.parser(s)
}

func (v param[T]) ParseAny(s string) (any, error) {
	return v.Parse(s)
}

func (v param[T]) Validators() Validators[T] {
	return v.validators
}

func (v param[T]) Validate(val T) error {
	return v.validators.Validate(val)
}

func (v param[T]) ValidateAny(val any) error {
	typed, ok := val.(T)
	if !ok {
		return fmt.Errorf("failed to cast %T to %T", val, typed)
	}

	return v.Validate(typed)
}

func (v param[T]) EncodeYAML(val T) (yaml.Node, error) {
	return v.yamlEncoder(val)
}

func (v param[T]) EncodeYAMLAny(val any) (yaml.Node, error) {
	typed, ok := val.(T)
	if !ok {
		return yaml.Node{}, fmt.Errorf("failed to cast %T to %T", val, typed)
	}

	return v.EncodeYAML(typed)
}

func (v param[T]) DecodeYAML(node yaml.Node) (T, error) {
	return v.yamlDecoder(node)
}

func (v param[T]) DecodeYAMLAny(node yaml.Node) (any, error) {
	return v.DecodeYAML(node)
}

func (v param[T]) ReflectType() reflect.Type {
	var t T
	return reflect.TypeOf(t)
}

func (v param[T]) DynamicType() (reflect.Type, bool) {
	return v.dynamicType, v.dynamicType != nil
}

func (v param[T]) IsDynamic() bool {
	return v.dynamicType != nil
}

func (v param[T]) EnumValues() []T {
	return v.enumValues
}

func comparatorReflect[T any](a, b T) bool {
	return reflect.ValueOf(a).Equal(reflect.ValueOf(b))
}

func stringifierSimple[T any](val T) string {
	return fmt.Sprintf("%v", val)
}

func WithStringifier[T any](stringifier func(T) string) Option[T] {
	return func(v *param[T]) error {
		v.stringifier = stringifier
		return nil
	}
}

func WithParser[T any](parser func(string) (T, error)) Option[T] {
	return func(v *param[T]) error {
		v.parser = parser
		return nil
	}
}

func WithEnumValues[T comparable](enumValues ...T) Option[T] {
	return func(v *param[T]) error {
		v.enumValues = enumValues
		return Validate(validatorOneOf[T]{
			comparator:  v.comparator,
			stringifier: v.stringifier,
			enumValues:  v.enumValues,
		})(v)
	}
}

func (v param[T]) Doc() string {
	return v.doc
}

// func (p param) ValidateDoc() string {
// 	return p.validateDoc
// }

// type typedParam[T any] struct {
// 	Value[T]
// 	param
// 	newDefault func() T
// }

func (v param[T]) NewDefault() T {
	return v.newDefault()
}

func (v param[T]) NewDefaultAny() any {
	return v.NewDefault()
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

func parserSlice[T any](elementParser func(string) (T, error)) func(string) ([]T, error) {
	return func(s string) ([]T, error) {
		parts, err := csv.NewReader(strings.NewReader(s)).Read()
		if err != nil {
			return nil, err
		}
		var result []T
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

func stringifierSlice[T any](elementStringifier func(T) string) func([]T) string {
	return func(sl []T) string {
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

// func validatorSlice[T any](elementValidator Validator[T]) Validator[[]T] {
// 	return validatorFunc[[]T](func(val []T) bool {
// 		for _, item := range val {
// 			if !elementValidator.Evaluate(item) {
// 				return false
// 			}
// 		}
// 		return true
// 	})
// }

type validatorSlice[T any] struct {
	elementValidator Validator[T]
}

func (v validatorSlice[T]) Doc() string {
	return "each element " + v.elementValidator.Doc()
}

func (v validatorSlice[T]) Evaluate(val []T) bool {
	return !slice.Some(val, func(item T) bool {
		return !v.elementValidator.Evaluate(item)
	})
}

type validatorDynamic[T any] struct {
	elementValidator Validator[T]
}

func (v validatorDynamic[T]) Doc() string {
	return "each element " + v.elementValidator.Doc()
}

func (v validatorDynamic[T]) Evaluate(val *atomic.Value[T]) bool {
	return v.elementValidator.Evaluate(val.Get())
}

func WithNewDefault[T any](newDefault func() T) Option[T] {
	return func(v *param[T]) error {
		v.newDefault = newDefault

		return nil
	}
}

func WithDefault[T comparable](defaultValue T) Option[T] {
	return WithNewDefault(func() T {
		return defaultValue
	})
}

type validatorRequired[T comparable] struct{}

func (v validatorRequired[T]) Doc() string {
	return "required"
}

func (v validatorRequired[T]) Evaluate(val T) bool {
	var zero T
	return val != zero
}

func WithRequired[T comparable]() Option[T] {
	return Validate(validatorRequired[T]{})
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type validatorGreaterThan[T number] struct {
	value T
}

func (v validatorGreaterThan[T]) Doc() string {
	return fmt.Sprintf("must be greater than %v", v.value)
}

func (v validatorGreaterThan[T]) Evaluate(val T) bool {
	return val > v.value
}

func WithGreaterThan[T number](min T) Option[T] {
	return Validate(validatorGreaterThan[T]{value: min})
}

type validatorLessThan[T number] struct {
	value T
}

func (v validatorLessThan[T]) Doc() string {
	return fmt.Sprintf("must be less than %v", v.value)
}

func (v validatorLessThan[T]) Evaluate(val T) bool {
	return val < v.value
}

func WithLessThan[T number](max T) Option[T] {
	return Validate(validatorLessThan[T]{value: max})
}

type validatorMin[T number] struct {
	value T
}

func (v validatorMin[T]) Doc() string {
	return fmt.Sprintf("must be at least %v", v.value)
}

func (v validatorMin[T]) Evaluate(val T) bool {
	return val >= v.value
}

func WithMin[T number](min T) Option[T] {
	return Validate(validatorMin[T]{value: min})
}

type validatorMax[T number] struct {
	value T
}

func (v validatorMax[T]) Doc() string {
	return fmt.Sprintf("must be at most %v", v.value)
}

func (v validatorMax[T]) Evaluate(val T) bool {
	return val <= v.value
}

func WithMax[T number](max T) Option[T] {
	return Validate(validatorMax[T]{value: max})
}
