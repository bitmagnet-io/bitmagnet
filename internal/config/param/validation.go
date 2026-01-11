package param

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type Validator[T any] interface {
	Doc() string
	Evaluate(T) error
}

var ErrInvalid = errors.New("validation failed")

func Validate[T any](validators ...Validator[T]) Option[T] {
	return func(p *param[T]) error {
		p.validators = append(p.validators, validators...)

		return nil
	}
}

type Validators[T any] []Validator[T]

func (v Validators[T]) Doc() string {
	return strings.Join(slice.Map(v, func(validator Validator[T]) string {
		return validator.Doc()
	}), "; ")
}

func (v Validators[T]) Evaluate(val T) error {
	return errors.Join(slice.Map(v, func(validator Validator[T]) error {
		return validator.Evaluate(val)
	})...)
}

func (v Validators[T]) Validate(val T) error {
	var errs []error
	for _, validator := range v {
		if err := validator.Evaluate(val); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%w: %s", ErrInvalid, errors.Join(errs...))
	}
	return nil
}

type minLengthValidator[T ~string] struct {
	min int
}

func (v minLengthValidator[T]) Doc() string {
	return fmt.Sprintf("must have minimum length %d", v.min)
}

func (v minLengthValidator[T]) Evaluate(val T) error {
	if len(val) >= v.min {
		return nil
	}
	return fmt.Errorf("length %d is less than minimum %d", len(val), v.min)
}

func MinLength[T ~string](min int) Option[T] {
	return Options(
		JSONSchemaOption[T](
			json_schema.MinLength(min),
			json_schema.Required(json_schema.RequiredBool(true)),
		),
		Validate(minLengthValidator[T]{min: min}),
	)
}

type maxLengthValidator[T ~string] struct {
	max int
}

func (v maxLengthValidator[T]) Doc() string {
	return fmt.Sprintf("must have maximum length %d", v.max)
}

func (v maxLengthValidator[T]) Evaluate(val T) error {
	if len(val) <= v.max {
		return nil
	}
	return fmt.Errorf("length %d is greater than maximum %d", len(val), v.max)
}

func MaxLength[T ~string](max int) Option[T] {
	return Options(
		JSONSchemaOption[T](json_schema.MaxLength(max)),
		Validate(maxLengthValidator[T]{max: max}),
	)
}

type validatorMinItems[E any, T ~[]E] struct {
	min int
}

func (v validatorMinItems[E, T]) Doc() string {
	return fmt.Sprintf("must have minimum items %d", v.min)
}

func (v validatorMinItems[E, T]) Evaluate(val T) error {
	if len(val) >= v.min {
		return nil
	}
	return fmt.Errorf("items count %d is less than minimum %d", len(val), v.min)
}

func MinItems[E any, T ~[]E](min int) Option[T] {
	return Options(
		JSONSchemaOption[T](json_schema.MinItems(min)),
		Validate(validatorMinItems[E, T]{min: min}),
	)
}

type validatorMaxItems[E any, T ~[]E] struct {
	max int
}

func (v validatorMaxItems[E, T]) Doc() string {
	return fmt.Sprintf("must have maximum items %d", v.max)
}

func (v validatorMaxItems[E, T]) Evaluate(val T) error {
	if len(val) <= v.max {
		return nil
	}
	return fmt.Errorf("items count %d is greater than maximum %d", len(val), v.max)
}

func MaxItems[E any, T ~[]E](max int) Option[T] {
	return Options(
		JSONSchemaOption[T](json_schema.MaxItems(max)),
		Validate(validatorMaxItems[E, T]{max: max}),
	)
}

type validatorOneOf[T any] struct {
	comparator  func(T, T) bool
	stringifier func(T) string
	enumValues  []T
}

func (v validatorOneOf[T]) Doc() string {
	return "must be one of: " + strings.Join(slice.Map(v.enumValues, v.stringifier), ", ")
}

func (v validatorOneOf[T]) Evaluate(val T) error {
	for _, enumVal := range v.enumValues {
		if v.comparator(val, enumVal) {
			return nil
		}
	}
	return fmt.Errorf("value %v is not one of the allowed values", val)
}

type validatorRegex[T ~string] struct {
	*regexp.Regexp
}

func (v validatorRegex[T]) Doc() string {
	return "must match pattern: " + v.Regexp.String()
}

func (v validatorRegex[T]) Evaluate(val T) error {
	if v.MatchString(string(val)) {
		return nil
	}
	return fmt.Errorf("value %q does not match pattern %q", val, v.Regexp.String())
}

func Regex[T ~string](regex *regexp.Regexp) Option[T] {
	return Options(
		JSONSchemaOption[T](json_schema.Pattern(regex)),
		Validate(validatorRegex[T]{Regexp: regex}),
	)
}

type validatorSlice[E any, T ~[]E] struct {
	elementValidator Validator[E]
}

func (v validatorSlice[E, T]) Doc() string {
	return "each element " + v.elementValidator.Doc()
}

func (v validatorSlice[E, T]) Evaluate(val T) error {
	return errors.Join(slice.Map(val, func(element E) error {
		return v.elementValidator.Evaluate(element)
	})...)
}

type validatorDynamic[T any] struct {
	elementValidator Validator[T]
}

func (v validatorDynamic[T]) Doc() string {
	return "each element " + v.elementValidator.Doc()
}

func (v validatorDynamic[T]) Evaluate(val *atomic.Value[T]) error {
	return v.elementValidator.Evaluate(val.Get())
}

type validatorRequired[T comparable] struct{}

func (v validatorRequired[T]) Doc() string {
	return "required"
}

func (v validatorRequired[T]) Evaluate(val T) error {
	var zero T
	if val != zero {
		return nil
	}
	return fmt.Errorf("value is required")
}

func Required[T comparable]() Option[T] {
	return Options(
		JSONSchemaOption[T](
			json_schema.Required(json_schema.RequiredBool(true)),
		),
		Validate(validatorRequired[T]{}),
	)
}

type validatorGreaterThan[T number] struct {
	value T
}

func (v validatorGreaterThan[T]) Doc() string {
	return fmt.Sprintf("must be greater than %v", v.value)
}

func (v validatorGreaterThan[T]) Evaluate(val T) error {
	if val > v.value {
		return nil
	}
	return fmt.Errorf("value %v is not greater than %v", val, v.value)
}

func GreaterThan[T number](min T) Option[T] {
	return Options(
		JSONSchemaOption[T](json_schema.ExclusiveMinimum(float64(min))),
		Validate(validatorGreaterThan[T]{value: min}),
	)
}

type validatorLessThan[T number] struct {
	value T
}

func (v validatorLessThan[T]) Doc() string {
	return fmt.Sprintf("must be less than %v", v.value)
}

func (v validatorLessThan[T]) Evaluate(val T) error {
	if val < v.value {
		return nil
	}
	return fmt.Errorf("value %v is not less than %v", val, v.value)
}

func LessThan[T number](max T) Option[T] {
	return Validate(validatorLessThan[T]{value: max})
}

type validatorMin[T number] struct {
	value T
}

func (v validatorMin[T]) Doc() string {
	return fmt.Sprintf("must be at least %v", v.value)
}

func (v validatorMin[T]) Evaluate(val T) error {
	if val >= v.value {
		return nil
	}
	return fmt.Errorf("value %v is not at least %v", val, v.value)
}

func Min[T number](min T) Option[T] {
	return Validate(validatorMin[T]{value: min})
}

type validatorMax[T number] struct {
	value T
}

func (v validatorMax[T]) Doc() string {
	return fmt.Sprintf("must be at most %v", v.value)
}

func (v validatorMax[T]) Evaluate(val T) error {
	if val <= v.value {
		return nil
	}
	return fmt.Errorf("value %v is not at most %v", val, v.value)
}

func Max[T number](max T) Option[T] {
	return Validate(validatorMax[T]{value: max})
}
