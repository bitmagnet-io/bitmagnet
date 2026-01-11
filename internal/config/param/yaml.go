package param

import (
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/go-viper/mapstructure/v2"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"
)

func yamlEncoder[T any](v T) (yaml.Node, error) {
	var node yaml.Node
	err := node.Encode(v)
	return node, err
}

func yamlDecoder[T any](node yaml.Node) (T, error) {
	var v T
	err := node.Decode(&v)
	return v, err
}

func yamlEncoderStringifier[T any](stringify func(T) string) func(v T) (yaml.Node, error) {
	return func(v T) (yaml.Node, error) {
		return yamlEncoder(stringify(v))
	}
}

func yamlDecoderParser[T any](parse func(string) (T, error)) func(yaml.Node) (T, error) {
	return func(node yaml.Node) (T, error) {
		str, err := yamlDecoder[string](node)
		if err != nil {
			var zero T
			return zero, err
		}

		return parse(str)
	}
}

func yamlEncoderSlice[E any, T ~[]E](elementEncoder func(E) (yaml.Node, error)) func(T) (yaml.Node, error) {
	return func(slice T) (yaml.Node, error) {
		node := yaml.Node{
			Kind: yaml.SequenceNode,
		}
		for _, elem := range slice {
			encodedElem, err := elementEncoder(elem)
			if err != nil {
				return yaml.Node{}, err
			}
			node.Content = append(node.Content, &encodedElem)
		}
		return node, nil
	}
}

func yamlDecoderSlice[E any, T ~[]E](elementDecoder func(yaml.Node) (E, error)) func(yaml.Node) (T, error) {
	return func(node yaml.Node) (T, error) {
		if node.Kind != yaml.SequenceNode {
			return nil, fmt.Errorf("expected sequence node, got %v", node.Kind)
		}
		var result T
		for _, item := range node.Content {
			decodedItem, err := elementDecoder(*item)
			if err != nil {
				return nil, err
			}
			result = append(result, decodedItem)
		}
		return result, nil
	}
}

func yamlEncoderDynamic[T any](elementEncoder func(T) (yaml.Node, error)) func(*atomic.Value[T]) (yaml.Node, error) {
	return func(value *atomic.Value[T]) (yaml.Node, error) {
		return elementEncoder(value.Get())
	}
}

func yamlDecoderDynamic[T any](elementDecoder func(yaml.Node) (T, error)) func(yaml.Node) (*atomic.Value[T], error) {
	return func(node yaml.Node) (*atomic.Value[T], error) {
		decodedItem, err := elementDecoder(node)
		if err != nil {
			return nil, err
		}

		return atomic.NewValue(decodedItem), nil
	}
}

func yamlDecoderMapstructure[T any]() func(node yaml.Node) (T, error) {
	return func(node yaml.Node) (T, error) {
		var value T

		var mapAny map[string]any
		err := node.Decode(&mapAny)
		if err != nil {
			return value, err
		}

		decoder, err := newDecoder(&value)

		if err != nil {
			return value, fmt.Errorf("failed to create decoder for type %T: %w", value, err)
		}

		if err := decoder.Decode(mapAny); err != nil {
			return value, fmt.Errorf("failed to decode value %v to type %T: %w", mapAny, value, err)
		}

		return value, nil
	}
}

func newDecoder(result any) (*mapstructure.Decoder, error) {
	return mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: result,
		MatchName: func(mapKey, fieldName string) bool {
			return mapKey == strcase.ToSnake(fieldName)
		},
		WeaklyTypedInput: true,
		Squash:           true,
		ErrorUnused:      true,
	})
}
