package json_schema

import (
	"errors"
	"regexp"

	"github.com/bitmagnet-io/bitmagnet/internal/ecma262"
)

// Option defines a functional option for configuring a Schema
type Option func(*JSONSchema) error

func Options(options ...Option) Option {
	return func(s *JSONSchema) error {
		for _, option := range options {
			err := option(s)

			if err != nil {
				return err
			}
		}

		return nil
	}
}

// Default sets the default value for the schema
func Default(value JSONValue) Option {
	return func(s *JSONSchema) error {
		if s.Default != nil {
			return errors.New("cannot overwrite existing default")
		}

		s.Default = &value
		return nil
	}
}

// Enum sets the enum values for the schema
func Enum(values ...JSONValue) Option {
	return func(s *JSONSchema) error {
		if len(values) == 0 {
			return errors.New("enum must have at least 1 value")
		}

		if s.Enum != nil {
			return errors.New("cannot overwrite existing enum")
		}

		s.Enum = values
		return nil
	}
}

// Pattern sets the pattern validation for the schema
func Pattern(re *regexp.Regexp) Option {
	return func(s *JSONSchema) error {
		if s.Pattern != nil {
			return errors.New("cannot overwrite existing pattern")
		}

		if s.Type != TypeString {
			return errors.New("pattern only appies to strings")
		}

		if err := ecma262.RegexpCompatibilityError(re); err != nil {
			return err
		}

		pattern := re.String()
		s.Pattern = &pattern

		return nil
	}
}

// MultipleOf sets the multipleOf validation for the schema
func MultipleOf(value float64) Option {
	return func(s *JSONSchema) error {
		if s.MultipleOf != nil {
			return errors.New("cannot overwrite existing multipleOf")
		}

		if !s.Type.IsNumeric() {
			return errors.New("multipleOf only applies to numeric types")
		}

		s.MultipleOf = &value
		return nil
	}
}

// Maximum sets the maximum validation for the schema
func Maximum(value float64) Option {
	return func(s *JSONSchema) error {
		if s.Maximum != nil && value > *s.Maximum {
			return errors.New("cannot overwrite existing maximum")
		}

		if !s.Type.IsNumeric() {
			return errors.New("maximum only applies to numeric types")
		}

		s.Maximum = &value
		return nil
	}
}

// ExclusiveMaximum sets the exclusiveMaximum validation for the schema
func ExclusiveMaximum(value float64) Option {
	return func(s *JSONSchema) error {
		if s.ExclusiveMaximum != nil && value > *s.ExclusiveMaximum {
			return errors.New("cannot overwrite existing exclusiveMaximum")
		}

		if !s.Type.IsNumeric() {
			return errors.New("exclusiveMaximum only applies to numeric types")
		}

		s.ExclusiveMaximum = &value
		return nil
	}
}

// Minimum sets the minimum validation for the schema
func Minimum(value float64) Option {
	return func(s *JSONSchema) error {
		if s.Minimum != nil && value < *s.Minimum {
			return errors.New("cannot overwrite existing minimum")
		}

		if !s.Type.IsNumeric() {
			return errors.New("minimum only applies to numeric types")
		}

		s.Minimum = &value
		return nil
	}
}

// ExclusiveMinimum sets the exclusiveMinimum validation for the schema
func ExclusiveMinimum(value float64) Option {
	return func(s *JSONSchema) error {
		if s.ExclusiveMaximum != nil && value < *s.ExclusiveMaximum {
			return errors.New("cannot overwrite existing exclusiveMinimum")
		}

		if !s.Type.IsNumeric() {
			return errors.New("exclusiveMinimum only applies to numeric types")
		}

		s.ExclusiveMinimum = &value
		return nil
	}
}

// MaxLength sets the maxLength validation for the schema
func MaxLength(value int) Option {
	return func(s *JSONSchema) error {
		if s.MaxLength != nil && value > *s.MaxLength {
			return errors.New("cannot overwrite existing maxLength")
		}

		if s.Type != TypeString {
			return errors.New("maxlength only applies to strings")
		}

		s.MaxLength = &value
		return nil
	}
}

// MinLength sets the minLength validation for the schema
func MinLength(value int) Option {
	return func(s *JSONSchema) error {
		if s.MinLength != nil && value < *s.MinLength {
			return errors.New("cannot overwrite existing minLength")
		}

		if s.Type != TypeString {
			return errors.New("minlength only applies to strings")
		}

		s.MinLength = &value
		return nil
	}
}

// MinItems sets the minItems validation for the schema
func MinItems(value int) Option {
	return func(s *JSONSchema) error {
		if s.MinItems != nil && value < *s.MinItems {
			return errors.New("cannot overwrite existing minItems")
		}

		if s.Type != TypeArray {
			return errors.New("minItems only applies to arrays")
		}

		s.MinItems = &value
		return nil
	}
}

// MaxItems sets the maxItems validation for the schema
func MaxItems(value int) Option {
	return func(s *JSONSchema) error {
		if s.MaxItems != nil && value > *s.MaxItems {
			return errors.New("cannot overwrite existing maxItems")
		}

		if s.Type != TypeArray {
			return errors.New("maxItems only applies to arrays")
		}

		s.MaxItems = &value
		return nil
	}
}

// UniqueItems sets the uniqueItems validation for the schema
func UniqueItems(value bool) Option {
	return func(s *JSONSchema) error {
		if s.UniqueItems != nil && value != *s.UniqueItems {
			return errors.New("cannot overwrite existing uniqueItems")
		}

		if s.Type != TypeArray {
			return errors.New("uniqueItems only applies to arrays")
		}

		s.UniqueItems = &value
		return nil
	}
}

// Required sets the required validation for the schema
func Required(value bool) Option {
	return func(s *JSONSchema) error {
		if s.Required != nil && value != *s.Required {
			return errors.New("cannot overwrite existing required")
		}

		s.Required = &value
		return nil
	}
}

// Items sets the items schema for the schema
func Items(schema JSONSchema) Option {
	return func(s *JSONSchema) error {
		if s.Items != nil {
			return errors.New("cannot overwrite existing items")
		}

		if s.Type != TypeArray {
			return errors.New("items only applies to arrays")
		}

		s.Items = &schema
		return nil
	}
}
