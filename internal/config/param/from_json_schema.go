package param

import (
	"errors"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

func JSONSchemaDecoder[T any](schema json_schema.JSONSchema) Option[T] {
	return func(p *param[T]) error {
		compiled, err := gojsonschema.NewSchema(gojsonschema.NewGoLoader(schema))
		if err != nil {
			return err
		}

		options := []Option[T]{
			JSONSchema[T](schema),
			jsonSchemaDecoder[T](compiled),
		}

		if schema.Description != nil {
			options = append(options, Description[T](*schema.Description))
		}

		if schema.Default != nil {
			var value T

			def := yaml.Node(*schema.Default)
			if err := def.Decode(&value); err != nil {
				return err
			}

			options = append(options, Default(value))
		}

		return Options(options...)(p)
	}
}

func jsonSchemaDecoder[T any](schema *gojsonschema.Schema) Option[T] {
	return func(p *param[T]) error {
		return YAMLDecoder[T](func(node yaml.Node) (T, error) {
			var raw any
			if err := node.Decode(&raw); err != nil {
				var zero T
				return zero, err
			}

			if err := jsonSchemaValidate(schema, raw); err != nil {
				var zero T
				return zero, err
			}

			var value T

			decoder, err := newDecoder(&value)
			if err != nil {
				var zero T
				return zero, err
			}

			if err := decoder.Decode(raw); err != nil {
				var zero T
				return zero, err
			}

			return value, nil
		})(p)
	}
}

func jsonSchemaValidate(schema *gojsonschema.Schema, val any) error {
	ld := gojsonschema.NewGoLoader(val)

	result, err := schema.Validate(ld)
	if err != nil {
		return err
	}

	if result.Valid() {
		return nil
	}

	var errs []error
	for _, desc := range result.Errors() {
		errs = append(errs, fmt.Errorf("%s", desc.String()))
	}

	return fmt.Errorf("schema validation failed: %w", errors.Join(errs...))
}
