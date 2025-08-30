package json_schema

import (
	"encoding/json"
	"io"
)

type JSONValue struct {
	Value any
}

func (v JSONValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Value)
}

func (v *JSONValue) UnmarshalJSON(data []byte) error {
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	v.Value = value

	return nil
}

func (v JSONValue) MarshalGQL(w io.Writer) {
	bytes, err := json.Marshal(v)
	if err == nil {
		_, _ = w.Write(bytes)
	}
}

func (v *JSONValue) UnmarshalGQL(gql any) error {
	if num, ok := gql.(json.Number); ok {
		flt, err := num.Float64()
		if err != nil {
			return err
		}

		v.Value = flt

		return nil
	}

	v.Value = gql

	return nil
}
