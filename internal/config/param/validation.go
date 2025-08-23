package param

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type Validator[T any] interface {
	Doc() string
	Evaluate(T) bool
}

func Validate[T any](validators ...Validator[T]) Option[T] {
	return func(v *param[T]) error {
		v.validators = append(v.validators, validators...)

		return nil
	}
}

type Validators[T any] []Validator[T]

func (v Validators[T]) Doc() string {
	return strings.Join(slice.Map(v, func(validator Validator[T]) string {
		return validator.Doc()
	}), "; ")
}

func (v Validators[T]) Evaluate(val T) bool {
	return !slice.Some(v, func(validator Validator[T]) bool {
		return !validator.Evaluate(val)
	})
}

var ErrInvalid = errors.New("validation failed")

func (v Validators[T]) Validate(val T) error {
	var errs []string
	for _, validator := range v {
		if !validator.Evaluate(val) {
			errs = append(errs, validator.Doc())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%w: %s", ErrInvalid, strings.Join(errs, "; "))
	}
	return nil
}

type lengthable[K comparable, V any] interface {
	~[]V | ~string | ~map[K]V
}

type minLengthValidator[K comparable, V any, T lengthable[K, V]] struct {
	min int
}

func (v minLengthValidator[K, V, T]) Doc() string {
	return fmt.Sprintf("must have minimum length %d", v.min)
}

func (v minLengthValidator[K, V, T]) Evaluate(val T) bool {
	return len(val) >= v.min
}

type maxLengthValidator[K comparable, V any, T lengthable[K, V]] struct {
	max int
}

func (v maxLengthValidator[K, V, T]) Doc() string {
	return fmt.Sprintf("must have maximum length %d", v.max)
}

func (v maxLengthValidator[K, V, T]) Evaluate(val T) bool {
	return len(val) <= v.max
}

func MinLength[K comparable, V any, T lengthable[K, V]](min int) Option[T] {
	return Validate(minLengthValidator[K, V, T]{min: min})
}

func MaxLength[K comparable, V any, T lengthable[K, V]](max int) Option[T] {
	return Validate(maxLengthValidator[K, V, T]{max: max})
}

type validatorOneOf[T any] struct {
	comparator  func(T, T) bool
	stringifier func(T) string
	enumValues  []T
}

func (v validatorOneOf[T]) Doc() string {
	return "must be one of: " + strings.Join(slice.Map(v.enumValues, v.stringifier), ", ")
}

func (v validatorOneOf[T]) Evaluate(val T) bool {
	for _, enumVal := range v.enumValues {
		if v.comparator(val, enumVal) {
			return true
		}
	}
	return false
}

type validatorRegex[T ~string] struct {
	*regexp.Regexp
}

func (v validatorRegex[T]) Doc() string {
	return "must match pattern: " + v.Regexp.String()
}

func (v validatorRegex[T]) Evaluate(val T) bool {
	return v.MatchString(string(val))
}

func WithValidatorRegex[T ~string](regex *regexp.Regexp) Option[T] {
	return Validate(validatorRegex[T]{Regexp: regex})
}
