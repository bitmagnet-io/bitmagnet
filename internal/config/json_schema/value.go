package json_schema

import (
	"encoding/json"
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

// JSONValue aliases yaml.Node as this package is used with both YAML and JSON payloads,
// and it is a useful and expressive type.
type JSONValue yaml.Node

func NewValue(value any) (JSONValue, error) {
	if num, ok := value.(json.Number); ok {
		flt, err := num.Float64()
		if err != nil {
			return JSONValue{}, err
		}

		value = flt
	}

	var node yaml.Node
	if err := node.Encode(value); err != nil {
		return JSONValue{}, err
	}

	return JSONValue(node), nil
}

func MustNewValue(value any) JSONValue {
	jv, err := NewValue(value)
	if err != nil {
		panic(err)
	}

	return jv
}

func (v JSONValue) yamlNode() *yaml.Node {
	vn := yaml.Node(v)
	return &vn
}

func (v JSONValue) Raw() any {
	var raw any
	_ = v.yamlNode().Decode(&raw)

	return raw
}

func (v JSONValue) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(v.yamlNode())
}

func (v JSONValue) MarshalJSON() ([]byte, error) {
	var value any
	if err := v.yamlNode().Decode(&value); err != nil {
		return nil, err
	}

	return json.Marshal(value)
}

func (v *JSONValue) UnmarshalYAML(data []byte) error {
	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return err
	}

	*v = JSONValue(node)

	return nil
}

func (v *JSONValue) UnmarshalJSON(data []byte) error {
	if !json.Valid(data) {
		return errors.New("invalid json")
	}

	return v.UnmarshalYAML(data)
}

func (v JSONValue) MarshalGQL(w io.Writer) {
	bytes, err := json.Marshal(v)
	if err == nil {
		_, _ = w.Write(bytes)
	}
}

func (v *JSONValue) UnmarshalGQL(gql any) error {
	node, err := NewValue(gql)
	if err != nil {
		return err
	}

	*v = node

	return nil
}
