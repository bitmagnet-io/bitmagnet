// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package classifier

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	FlagTypeBool            FlagType = "bool"
	FlagTypeString          FlagType = "string"
	FlagTypeInt             FlagType = "int"
	FlagTypeStringList      FlagType = "string_list"
	FlagTypeContentTypeList FlagType = "content_type_list"
)

var ErrInvalidFlagType = fmt.Errorf("not a valid FlagType, try [%s]", strings.Join(_FlagTypeNames, ", "))

var _FlagTypeNames = []string{
	string(FlagTypeBool),
	string(FlagTypeString),
	string(FlagTypeInt),
	string(FlagTypeStringList),
	string(FlagTypeContentTypeList),
}

// FlagTypeNames returns a list of possible string values of FlagType.
func FlagTypeNames() []string {
	tmp := make([]string, len(_FlagTypeNames))
	copy(tmp, _FlagTypeNames)
	return tmp
}

// FlagTypeValues returns a list of the values for FlagType
func FlagTypeValues() []FlagType {
	return []FlagType{
		FlagTypeBool,
		FlagTypeString,
		FlagTypeInt,
		FlagTypeStringList,
		FlagTypeContentTypeList,
	}
}

// String implements the Stringer interface.
func (x FlagType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x FlagType) IsValid() bool {
	_, err := ParseFlagType(string(x))
	return err == nil
}

var _FlagTypeValue = map[string]FlagType{
	"bool":              FlagTypeBool,
	"string":            FlagTypeString,
	"int":               FlagTypeInt,
	"string_list":       FlagTypeStringList,
	"content_type_list": FlagTypeContentTypeList,
}

// ParseFlagType attempts to convert a string to a FlagType.
func ParseFlagType(name string) (FlagType, error) {
	if x, ok := _FlagTypeValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _FlagTypeValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return FlagType(""), fmt.Errorf("%s is %w", name, ErrInvalidFlagType)
}

// MarshalText implements the text marshaller method.
func (x FlagType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *FlagType) UnmarshalText(text []byte) error {
	tmp, err := ParseFlagType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

var errFlagTypeNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *FlagType) Scan(value interface{}) (err error) {
	if value == nil {
		*x = FlagType("")
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case string:
		*x, err = ParseFlagType(v)
	case []byte:
		*x, err = ParseFlagType(string(v))
	case FlagType:
		*x = v
	case *FlagType:
		if v == nil {
			return errFlagTypeNilPtr
		}
		*x = *v
	case *string:
		if v == nil {
			return errFlagTypeNilPtr
		}
		*x, err = ParseFlagType(*v)
	default:
		return errors.New("invalid type for FlagType")
	}

	return
}

// Value implements the driver Valuer interface.
func (x FlagType) Value() (driver.Value, error) {
	return x.String(), nil
}

type NullFlagType struct {
	FlagType FlagType
	Valid    bool
	Set      bool
}

func NewNullFlagType(val interface{}) (x NullFlagType) {
	err := x.Scan(val) // yes, we ignore this error, it will just be an invalid value.
	_ = err            // make any errcheck linters happy
	return
}

// Scan implements the Scanner interface.
func (x *NullFlagType) Scan(value interface{}) (err error) {
	if value == nil {
		x.FlagType, x.Valid = FlagType(""), false
		return
	}

	err = x.FlagType.Scan(value)
	x.Valid = (err == nil)
	return
}

// Value implements the driver Valuer interface.
func (x NullFlagType) Value() (driver.Value, error) {
	if !x.Valid {
		return nil, nil
	}
	return x.FlagType.String(), nil
}

// MarshalJSON correctly serializes a NullFlagType to JSON.
func (n NullFlagType) MarshalJSON() ([]byte, error) {
	const nullStr = "null"
	if n.Valid {
		return json.Marshal(n.FlagType)
	}
	return []byte(nullStr), nil
}

// UnmarshalJSON correctly deserializes a NullFlagType from JSON.
func (n *NullFlagType) UnmarshalJSON(b []byte) error {
	n.Set = true
	var x interface{}
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	err = n.Scan(x)
	return err
}
