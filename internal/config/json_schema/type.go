package json_schema

import (
	"encoding/json"
	"fmt"
)

type Type string

const (
	TypeString  Type = "string"
	TypeNumber  Type = "number"
	TypeInteger Type = "integer"
	TypeBoolean Type = "boolean"
	TypeArray   Type = "array"
)

func (t Type) String() string {
	return string(t)
}

func (t Type) Valid() error {
	switch t {
	case TypeString, TypeNumber, TypeInteger, TypeBoolean, TypeArray:
		return nil
	default:
		return fmt.Errorf("invalid Type value: %q", t)
	}
}

func (t Type) IsNumeric() bool {
	return t == TypeNumber || t == TypeInteger
}

func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Type) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal Type: %w", err)
	}

	*t = Type(s)

	return t.Valid()
}
