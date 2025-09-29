package json_schema

import (
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/ecma262"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
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

func Metaschema(ms string) Option {
	return func(s *JSONSchema) error {
		if s.Schema != nil {
			return errors.New("cannot overwrite existing $schema")
		}

		s.Schema = &ms

		return nil
	}
}

func MetaschemaDraft7() Option {
	return Metaschema("http://json-schema.org/draft-07/schema#")
}

func ID(id string) Option {
	return func(s *JSONSchema) error {
		if s.ID != nil {
			return errors.New("cannot overwrite existing $id")
		}

		s.ID = &id

		return nil
	}
}

func Typed(tp Type) Option {
	return func(s *JSONSchema) error {
		if s.Type != nil {
			return errors.New("cannot overwrite existing type")
		}

		if s.Ref != nil {
			return errors.New("cannot set both ref and type")
		}

		s.Type = &tp

		return nil
	}
}

func Ref(ref string) Option {
	return func(s *JSONSchema) error {
		if s.Ref != nil {
			return errors.New("cannot overwrite existing ref")
		}

		if s.Type != nil {
			return errors.New("cannot set both type and ref")
		}

		s.Ref = &ref

		return nil
	}
}

func RefDefinition(name string) Option {
	return Ref("#/definitions/" + name)
}

func DescriptionIfNonEmpty(str string) Option {
	return func(s *JSONSchema) error {
		if str != "" {
			return Description(str)(s)
		}

		return nil
	}
}

func Description(str string) Option {
	return func(s *JSONSchema) error {
		if s.Description != nil {
			return errors.New("cannot overwrite existing description")
		}

		s.Description = &str

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

func Const(value JSONValue) Option {
	return func(s *JSONSchema) error {
		if s.Const != nil {
			return errors.New("cannot overwrite existing const")
		}

		s.Const = &value

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

		if !s.HasType(TypeString) {
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

		if !s.HasType(TypeString) {
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

		if !s.HasType(TypeString) {
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

		if !s.HasType(TypeArray) {
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

		if !s.HasType(TypeArray) {
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

		if !s.HasType(TypeArray) {
			return errors.New("uniqueItems only applies to arrays")
		}

		s.UniqueItems = &value
		return nil
	}
}

// Required sets the required validation for the schema
func Required(value RequiredParam) Option {
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

		if !s.HasType(TypeArray) {
			return errors.New("items only applies to arrays")
		}

		s.Items = &schema
		return nil
	}
}

func Properties(props map[string]JSONSchema) Option {
	return func(s *JSONSchema) error {
		if s.Properties != nil {
			return errors.New("cannot overwrite existing properties")
		}

		if !s.HasType(TypeObject) {
			return errors.New("properties only applies to objects")
		}

		s.Properties = props

		return nil
	}
}

func AdditionalProperties(value AdditionalPropertiesParam) Option {
	return func(s *JSONSchema) error {
		if !s.HasType(TypeObject) {
			return errors.New("additionalProperties only applies to objects")
		}

		s.AdditionalProperties = &value

		return nil
	}
}

func AdditionalPropertiesType(schema JSONSchema) Option {
	return AdditionalProperties(AdditionalPropertiesSchema(schema))
}

func AdditionalPropertiesTrue() Option {
	b := AdditionalPropertiesBool(true)
	return AdditionalProperties(&b)
}

func AdditionalPropertiesFalse() Option {
	b := AdditionalPropertiesBool(false)
	return AdditionalProperties(&b)
}

func OneOf(schemas ...JSONSchema) Option {
	return func(s *JSONSchema) error {
		if s.OneOf != nil {
			return errors.New("cannot overwrite existing oneOf")
		}

		s.OneOf = schemas

		return nil
	}
}

func Definition(name string, schema JSONSchema) Option {
	return func(s *JSONSchema) error {
		if s.Definitions == nil {
			s.Definitions = make(map[string]JSONSchema)
		}

		if _, ok := s.Definitions[name]; ok {
			return fmt.Errorf("cannot overwrite existing definition: %s", name)
		}

		s.Definitions[name] = schema

		return nil
	}
}

func Definitions(defs map[string]JSONSchema) Option {
	return Options(slice.Map(slices.Sorted(maps.Keys(defs)), func(k string) Option {
		return Definition(k, defs[k])
	})...)
}
