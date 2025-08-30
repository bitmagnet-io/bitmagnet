package param

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/ecma262"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"gopkg.in/yaml.v3"
)

type Option[T any] func(*param[T]) error

func Options[T any](opts ...Option[T]) Option[T] {
	return func(p *param[T]) error {
		for _, opt := range opts {
			if err := opt(p); err != nil {
				return err
			}
		}
		return nil
	}
}

func StringLiteral[T ~string]() Option[T] {
	return Options(
		Parser(func(s string) (T, error) {
			return T(s), nil
		}),
	)
}

func Slice[E any, T ~[]E](opts ...Option[E]) Option[T] {
	return func(p *param[T]) error {
		elem, err := New(append(opts, nested[E]())...)
		if err != nil {
			return err
		}

		return Options(
			JSONSchemaTypeOptions[T](json_schema.TypeArray, json_schema.Items(elem.JSONSchema())),
			Stringifier(stringifierSlice[E, T](elem.Stringify)),
			Parser(parserSlice[E, T](elem.Parse)),
			YAMLEncoder(yamlEncoderSlice[E, T](elem.EncodeYAML)),
			YAMLDecoder(yamlDecoderSlice[E, T](elem.DecodeYAML)),
			Validate(slice.Map(elem.Validators(), func(validator Validator[E]) Validator[T] {
				return validatorSlice[E, T]{elementValidator: validator}
			})...),
		)(p)
	}
}

func Dynamic[T any](opts ...Option[T]) Option[*atomic.Value[T]] {
	return func(p *param[*atomic.Value[T]]) error {
		elem, err := New(append(opts, nested[T]())...)
		if err != nil {
			return err
		}

		description := p.description
		if elemDoc := elem.Description(); elemDoc != "" {
			description = elemDoc
		}

		return Options(
			dynamicType[T](elem.ReflectType()),
			JSONSchema[*atomic.Value[T]](elem.JSONSchema()),
			Stringifier(stringifierDynamic(elem.Stringify)),
			Parser(parserDynamic(elem.Parse)),
			YAMLEncoder(yamlEncoderDynamic(elem.EncodeYAML)),
			YAMLDecoder(yamlDecoderDynamic(elem.DecodeYAML)),
			Validate(slice.Map(elem.Validators(), func(validator Validator[T]) Validator[*atomic.Value[T]] {
				return validatorDynamic[T]{elementValidator: validator}
			})...),
			NewDefault(func() *atomic.Value[T] {
				return atomic.NewValue(elem.NewDefault())
			}),
			Description[*atomic.Value[T]](description),
			func(p *param[*atomic.Value[T]]) error {
				p.comparator = func(a, b *atomic.Value[T]) bool {
					return elem.Equals(a.Get(), b.Get())
				}

				if enumValues := elem.EnumValues(); enumValues != nil {
					p.enumValues = slice.Map(enumValues, atomic.NewValue)
				}

				return nil
			},
		)(p)
	}
}

func nested[T any]() Option[T] {
	return func(p *param[T]) error {
		p.nested = true

		return nil
	}
}

func dynamicType[T any](reflectType reflect.Type) Option[*atomic.Value[T]] {
	return func(p *param[*atomic.Value[T]]) error {
		p.dynamicType = reflectType

		return nil
	}
}

func Mapstructure[T any]() Option[T] {
	return YAMLDecoder(yamlDecoderMapstructure[T]())
}

func YAMLEncoder[T any](encoder func(T) (yaml.Node, error)) Option[T] {
	return func(p *param[T]) error {
		p.yamlEncoder = encoder
		return nil
	}
}

func YAMLDecoder[T any](decoder func(yaml.Node) (T, error)) Option[T] {
	return func(p *param[T]) error {
		p.yamlDecoder = decoder
		return nil
	}
}

func Comparator[T any](compare func(T, T) bool) Option[T] {
	return func(p *param[T]) error {
		p.comparator = compare

		return nil
	}
}

func Description[T any](doc string) Option[T] {
	return func(p *param[T]) error {
		p.description = doc

		return nil
	}
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type float interface {
	~float32 | ~float64
}

func Int[T integer]() Option[T] {
	return JSONSchemaTypeOptions[T](json_schema.TypeInteger)
}

func Float[T float]() Option[T] {
	return JSONSchemaTypeOptions[T](json_schema.TypeNumber)
}

func Bool[T ~bool]() Option[T] {
	return JSONSchemaTypeOptions[T](json_schema.TypeBoolean)
}

var (
	regexDuration         = ecma262.MustCompile(`^([-+]?((\d+(\.\d*)?|\.\d+)(ns|us|\\u00b5s|\\u03bcs|ms|s|m|h))+|0)$`)
	regexPositiveDuration = ecma262.MustCompile(`^(\+?((\d+(\.\d*)?|\.\d+)(ns|us|\\u00b5s|\\u03bcs|ms|s|m|h))+)$`)
)

func Duration[T ~int64](positive bool) Option[T] {
	re := regexDuration
	if positive {
		re = regexPositiveDuration
	}

	return Options(
		JSONSchemaOption[T](json_schema.Pattern(re)),
		func(p *param[T]) error {
			if positive {
				return Validate(validatorGreaterThan[T]{value: 0})(p)
			}

			return nil
		},
		Stringifier(stringifyDuration[T]),
		Parser(parseDuration[T]),
		YAMLEncoder(yamlEncoderStringifier(stringifyDuration[T])),
		YAMLDecoder(yamlDecoderParser(parseDuration[T])),
	)
}

func stringifyDuration[T ~int64](value T) string {
	return time.Duration(value).String()
}

func parseDuration[T ~int64](str string) (T, error) {
	duration, err := time.ParseDuration(str)
	if err != nil {
		return 0, err
	}

	return T(duration), nil
}

func JSONSchema[T any](schema json_schema.JSONSchema) Option[T] {
	return func(p *param[T]) error {
		if !p.jsonSchema.IsBasicString() {
			return errors.New("cannot overwrite pre-existing schema")
		}

		p.jsonSchema = schema

		return nil
	}
}

func JSONSchemaTypeOptions[T any](schemaType json_schema.Type, options ...json_schema.Option) Option[T] {
	return func(p *param[T]) error {
		if schemaType != p.jsonSchema.Type && !p.jsonSchema.IsBasicString() {
			return errors.New("cannot overwrite pre-existing schema")
		}

		p.jsonSchema.Type = schemaType

		return JSONSchemaOption[T](options...)(p)
	}
}

func JSONSchemaOption[T any](options ...json_schema.Option) Option[T] {
	return func(p *param[T]) error {
		return json_schema.Options(options...)(&p.jsonSchema)
	}
}

func comparatorReflect[T any](a, b T) bool {
	return reflect.ValueOf(a).Equal(reflect.ValueOf(b))
}

func stringifierSimple[T any](val T) string {
	return fmt.Sprintf("%v", val)
}

func Stringifier[T any](stringifier func(T) string) Option[T] {
	return func(p *param[T]) error {
		p.stringifier = stringifier
		return nil
	}
}

func Parser[T any](parser func(string) (T, error)) Option[T] {
	return func(v *param[T]) error {
		v.parser = parser
		return nil
	}
}

func EnumValues[T comparable](enumValues ...T) Option[T] {
	return func(p *param[T]) error {
		p.enumValues = enumValues

		jsonValues := make([]json_schema.JSONValue, 0, len(enumValues))

		for _, val := range enumValues {
			enc, err := p.EncodeYAMLAnyAny(val)
			if err != nil {
				return err
			}
			jsonValues = append(jsonValues, json_schema.JSONValue{Value: enc})
		}

		return Options(
			JSONSchemaOption[T](json_schema.Enum(jsonValues...)),
			Validate(validatorOneOf[T]{
				comparator:  p.comparator,
				stringifier: p.stringifier,
				enumValues:  p.enumValues,
			}),
		)(p)
	}
}

func NewDefault[T any](newDefault func() T) Option[T] {
	return func(p *param[T]) error {
		p.newDefault = newDefault
		p.hasExplicitDefault = true

		return nil
	}
}

func Default[T comparable](defaultValue T) Option[T] {
	return NewDefault(func() T {
		return defaultValue
	})
}

func PortNumber[T ~uint16]() Option[T] {
	return Options(
		Int[T](),
		Min(T(1)),
		Max(T(65535)),
	)
}
