package configresolver

import (
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"os"
)

func NewFromYamlFile(path string, ignoreMissing bool, val *validator.Validate, options ...Option) (Resolver, error) {
	m := make(map[string]interface{})
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		if !ignoreMissing || !os.IsNotExist(readErr) {
			return nil, readErr
		}
	} else {
		parseErr := yaml.Unmarshal(data, &m)
		if parseErr != nil {
			return nil, parseErr
		}
	}
	return NewMap(m, val, append([]Option{WithKey(path)}, options...)...), nil
}
