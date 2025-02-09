package classifier

import (
	"github.com/go-viper/mapstructure/v2"

)

func newDecoder[T any](target *T) (*mapstructure.Decoder, error) {
	return mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: target,
		MatchName: func(mapKey, fieldName string) bool {
			return mapKey == strcase.ToSnake(fieldName)
		},
		ErrorUnused: true,
		TagName:     "json",
	})
}
