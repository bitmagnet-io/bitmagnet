package param

import (
	"reflect"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"gopkg.in/yaml.v3"
)

type Option[T any] func(*param[T]) error

func WithOptions[T any](opts ...Option[T]) Option[T] {
	return func(v *param[T]) error {
		for _, opt := range opts {
			if err := opt(v); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithStringLiteral[T ~string]() Option[T] {
	return WithOptions(
		WithParser(func(s string) (T, error) {
			return T(s), nil
		}),
	)
}

func WithSlice[T any](opts ...Option[T]) Option[[]T] {
	return func(p *param[[]T]) error {
		elem, err := New(opts...)
		if err != nil {
			return err
		}
		return WithOptions(
			WithStringifier(stringifierSlice(elem.Stringify)),
			WithParser(parserSlice(elem.Parse)),
			WithYAMLEncoder(yamlEncoderSlice(elem.EncodeYAML)),
			WithYAMLDecoder(yamlDecoderSlice(elem.DecodeYAML)),
			Validate(slice.Map(elem.Validators(), func(validator Validator[T]) Validator[[]T] {
				return validatorSlice[T]{elementValidator: validator}
			})...),
		)(p)
	}
}

func WithDynamic[T any](opts ...Option[T]) Option[*atomic.Value[T]] {
	return func(v *param[*atomic.Value[T]]) error {
		elem, err := New(opts...)
		if err != nil {
			return err
		}

		doc := v.doc
		if elemDoc := elem.Doc(); elemDoc != "" {
			doc = elemDoc
		}

		return WithOptions(
			withDynamicType[T](elem.ReflectType()),
			WithStringifier(stringifierDynamic(elem.Stringify)),
			WithParser(parserDynamic(elem.Parse)),
			WithYAMLEncoder(yamlEncoderDynamic(elem.EncodeYAML)),
			WithYAMLDecoder(yamlDecoderDynamic(elem.DecodeYAML)),
			Validate(slice.Map(elem.Validators(), func(validator Validator[T]) Validator[*atomic.Value[T]] {
				return validatorDynamic[T]{elementValidator: validator}
			})...),
			WithRequired[*atomic.Value[T]](),
			WithNewDefault(func() *atomic.Value[T] {
				return atomic.NewValue(elem.NewDefault())
			}),
			WithDoc[*atomic.Value[T]](doc),
		)(v)
	}
}

func withDynamicType[T any](reflectType reflect.Type) Option[*atomic.Value[T]] {
	return func(param *param[*atomic.Value[T]]) error {
		param.dynamicType = reflectType

		return nil
	}
}

func WithMapstructure[T any]() Option[T] {
	return WithYAMLDecoder(yamlDecoderMapstructure[T]())
}

func WithYAMLEncoder[T any](encoder func(T) (yaml.Node, error)) Option[T] {
	return func(p *param[T]) error {
		p.yamlEncoder = encoder
		return nil
	}
}

func WithYAMLDecoder[T any](decoder func(yaml.Node) (T, error)) Option[T] {
	return func(p *param[T]) error {
		p.yamlDecoder = decoder
		return nil
	}
}

func WithComparator[T any](compare func(T, T) bool) Option[T] {
	return func(p *param[T]) error {
		p.comparator = compare

		return nil
	}
}

func WithDoc[T any](doc string) Option[T] {
	return func(p *param[T]) error {
		p.doc = doc

		return nil
	}
}
