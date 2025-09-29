package param

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/json_schema"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

func Slice[E any, T ~[]E](opts ...Option[E]) Option[T] {
	return func(p *param[T]) error {
		valueParam, err := New(opts...)
		if err != nil {
			return err
		}

		return Options(
			JSONSchemaTypeOptions[T](json_schema.TypeArray, json_schema.Items(valueParam.JSONSchema())),
			Stringifier(stringifierSlice[E, T](valueParam.Stringify)),
			Parser(parserSlice[E, T](valueParam.Parse)),
			YAMLEncoder(yamlEncoderSlice[E, T](valueParam.EncodeYAML)),
			YAMLDecoder(yamlDecoderSlice[E, T](valueParam.DecodeYAML)),
			Validate(slice.Map(valueParam.Validators(), func(validator Validator[E]) Validator[T] {
				return validatorSlice[E, T]{elementValidator: validator}
			})...,
			))(p)
	}
}
