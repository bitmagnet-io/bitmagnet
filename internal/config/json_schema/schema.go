package json_schema

// JSONSchema represents a subset of the JSON schema spacification that can be used for configuration parameters.
type JSONSchema struct {
	Type             Type        `json:"type"`
	Default          *JSONValue  `json:"default,omitempty"`
	Enum             []JSONValue `json:"enum,omitempty"`
	Pattern          *string     `json:"pattern,omitempty"`
	MultipleOf       *float64    `json:"multipleOf,omitempty"`
	Maximum          *float64    `json:"maximum,omitempty"`
	ExclusiveMaximum *float64    `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64    `json:"minimum,omitempty"`
	ExclusiveMinimum *float64    `json:"exclusiveMinimum,omitempty"`
	MaxLength        *int        `json:"maxLength,omitempty"`
	MinLength        *int        `json:"minLength,omitempty"`
	MinItems         *int        `json:"minItems,omitempty"`
	MaxItems         *int        `json:"maxItems,omitempty"`
	UniqueItems      *bool       `json:"uniqueItems,omitempty"`
	Required         *bool       `json:"required,omitempty"`
	Items            *JSONSchema `json:"items,omitempty"`
}

// New creates a new Schema with the given type and optional configuration
func New(schemaType Type, options ...Option) (JSONSchema, error) {
	s := JSONSchema{
		Type: schemaType,
	}

	err := Options(options...)(&s)

	return s, err
}

func MustNew(schemaType Type, options ...Option) JSONSchema {
	schema, err := New(schemaType, options...)
	if err != nil {
		panic(err)
	}

	return schema
}

func (s JSONSchema) IsBasicString() bool {
	return s.Type == TypeString && !s.HasRules()
}

func (s JSONSchema) HasRules() bool {
	return s.Pattern != nil ||
		s.Enum != nil ||
		s.MultipleOf != nil ||
		s.Maximum != nil ||
		s.ExclusiveMaximum != nil ||
		s.Minimum != nil ||
		s.ExclusiveMinimum != nil ||
		s.MaxLength != nil ||
		s.MinLength != nil ||
		s.MinItems != nil ||
		s.MaxItems != nil ||
		s.UniqueItems != nil ||
		s.Required != nil ||
		s.Items != nil
}
