package cmd

import (
	"bytes"
	"encoding"
	"encoding/csv"
	"errors"
)

type CSV[T any] []T

func (c *CSV[T]) UnmarshalText(text []byte) error {
	csValues, err := csv.NewReader(bytes.NewBuffer(text)).Read()
	if err != nil {
		return err
	}

	values := make([]T, 0, len(csValues))

	for _, strValue := range csValues {
		value, err := c.parse(strValue)
		if err != nil {
			return err
		}
		values = append(values, value)
	}

	*c = values

	return nil
}

func (c *CSV[T]) parse(str string) (T, error) {
	var result T

	value := new(T)

	if unm, ok := any(value).(encoding.TextUnmarshaler); ok {
		err := unm.UnmarshalText([]byte(str))
		if err != nil {
			return result, err
		}
		return *value, nil
	}

	return result, errors.New("unmarshal failed")
}

type CSVStringSlice []string

func (c *CSVStringSlice) UnmarshalText(text []byte) error {
	csvs, err := csv.NewReader(bytes.NewBuffer(text)).Read()
	if err != nil {
		return err
	}

	*c = csvs

	return nil
}
