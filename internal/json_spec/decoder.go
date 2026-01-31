package json_spec

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/iancoleman/strcase"
)

type KeyMatcher = func(mapKey, fieldName string) bool

func KeyMatcherSnake(mapKey, fieldName string) bool {
	return mapKey == strcase.ToSnake(fieldName)
}

func KeyMatcherLowerCamel(mapKey, fieldName string) bool {
	return mapKey == strcase.ToLowerCamel(fieldName)
}

func NewDecoder[T any](keyMatcher KeyMatcher, target *T) (*mapstructure.Decoder, error) {
	return mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:      target,
		MatchName:   keyMatcher,
		ErrorUnused: true,
		TagName:     "json",
	})
}

func DecodeTo[T any](ctx ParseContext, target *T) error {
	decoder, decoderErr := NewDecoder(ctx.KeyMatcher, target)
	if decoderErr != nil {
		return ctx.Error(decoderErr)
	}

	return decoder.Decode(ctx.Source)
}

func Decode[T any](ctx ParseContext) (T, error) {
	var target T

	err := DecodeTo(ctx, &target)

	return target, err
}
