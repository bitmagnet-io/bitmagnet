package json_schema

import (
	"encoding/json"
	"errors"
	"slices"
)

// JSONSchema represents a subset of the JSON schema spacification.
type JSONSchema struct {
	Schema               *string                    `json:"$schema,omitempty" yaml:"$schema,omitempty"`
	ID                   *string                    `json:"$id,omitempty" yaml:"$id,omitempty"`
	Ref                  *string                    `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Type                 *Type                      `json:"type,omitempty" yaml:"type,omitempty"`
	Description          *string                    `json:"description,omitempty" yaml:"description,omitempty"`
	Default              *JSONValue                 `json:"default,omitempty" yaml:"default,omitempty"`
	Const                *JSONValue                 `json:"const,omitempty" yaml:"const,omitempty"`
	Enum                 []JSONValue                `json:"enum,omitempty" yaml:"enum,omitempty"`
	Pattern              *string                    `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MultipleOf           *float64                   `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	Maximum              *float64                   `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	ExclusiveMaximum     *float64                   `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	Minimum              *float64                   `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	ExclusiveMinimum     *float64                   `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`
	MaxLength            *int                       `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength            *int                       `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MinItems             *int                       `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	MaxItems             *int                       `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	UniqueItems          *bool                      `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	Required             *RequiredParam             `json:"required,omitempty" yaml:"required,omitempty"`
	Nullable             *bool                      `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	Items                *JSONSchema                `json:"items,omitempty" yaml:"items,omitempty"`
	Properties           map[string]JSONSchema      `json:"properties,omitempty" yaml:"properties,omitempty"`
	AdditionalProperties *AdditionalPropertiesParam `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	OneOf                []JSONSchema               `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	Definitions          map[string]JSONSchema      `json:"definitions,omitempty" yaml:"definitions,omitempty"`
}

func (s *JSONSchema) UnmarshalJSON(data []byte) error {
	type Alias JSONSchema // Prevent recursion
	aux := &struct {
		AdditionalProperties json.RawMessage `json:"additionalProperties,omitempty"`
		Required             json.RawMessage `json:"required,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux.AdditionalProperties) > 0 {
		ap, err := UnmarshalAdditionalProperties(aux.AdditionalProperties)
		if err != nil {
			return err
		}

		s.AdditionalProperties = &ap
	}

	if len(aux.Required) > 0 {
		required, err := UnmarshalRequired(aux.Required)
		if err != nil {
			return err
		}

		s.Required = &required
	}

	return nil
}

func (s JSONSchema) HasType(tps ...Type) bool {
	return s.Type != nil && slices.Contains(tps, *s.Type)
}

// New creates a new Schema with the given configuration
func New(options ...Option) (JSONSchema, error) {
	s := JSONSchema{}

	err := Options(options...)(&s)

	return s, err
}

func MustNew(options ...Option) JSONSchema {
	schema, err := New(options...)
	if err != nil {
		panic(err)
	}

	return schema
}

func (s JSONSchema) IsBasicString() bool {
	return s.HasType(TypeString) && !s.HasRules()
}

func (s JSONSchema) HasRules() bool {
	return s.Const != nil ||
		s.Enum != nil ||
		s.ExclusiveMaximum != nil ||
		s.ExclusiveMinimum != nil ||
		s.Items != nil ||
		s.Maximum != nil ||
		s.MaxItems != nil ||
		s.MaxLength != nil ||
		s.Minimum != nil ||
		s.MinItems != nil ||
		s.MinLength != nil ||
		s.MultipleOf != nil ||
		s.OneOf != nil ||
		s.Pattern != nil ||
		s.Ref != nil ||
		s.Required != nil ||
		s.UniqueItems != nil
}

type RequiredParam interface {
	required()
}

type RequiredBool bool

func (r *RequiredBool) UnmarshalJSON(data []byte) error {
	var b bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}

	*r = RequiredBool(b)

	return nil
}

func (RequiredBool) required() {}

type RequiredFields []string

func (r *RequiredFields) UnmarshalJSON(data []byte) error {
	var fields []string

	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	*r = RequiredFields(fields)

	return nil
}

func (RequiredFields) required() {}

func UnmarshalRequired(data []byte) (RequiredParam, error) {
	var errs []error

	var b RequiredBool

	err := json.Unmarshal(data, &b)
	if err == nil {
		return b, nil
	}

	errs = append(errs, err)

	var fields RequiredFields

	err = json.Unmarshal(data, &fields)
	if err == nil {
		return fields, nil
	}

	errs = append(errs, err)

	return nil, errors.Join(errs...)
}

type AdditionalPropertiesParam interface {
	additionalProperties()
}

type AdditionalPropertiesBool bool

func (AdditionalPropertiesBool) additionalProperties() {}

func (p *AdditionalPropertiesBool) UnmarshalJSON(data []byte) error {
	var b bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}

	*p = AdditionalPropertiesBool(b)

	return nil
}

type AdditionalPropertiesSchema JSONSchema

func (p *AdditionalPropertiesSchema) UnmarshalJSON(data []byte) error {
	var s JSONSchema
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*p = AdditionalPropertiesSchema(s)

	return nil
}

func (AdditionalPropertiesSchema) additionalProperties() {}

func UnmarshalAdditionalProperties(data []byte) (AdditionalPropertiesParam, error) {
	var errs []error

	var b AdditionalPropertiesBool

	err := json.Unmarshal(data, &b)
	if err == nil {
		return b, nil
	}

	errs = append(errs, err)

	var schema AdditionalPropertiesSchema

	err = json.Unmarshal(data, &schema)
	if err == nil {
		return schema, nil
	}

	errs = append(errs, err)

	return nil, errors.Join(errs...)
}
