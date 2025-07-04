package configfx

import (
	"errors"

	"github.com/bitmagnet-io/bitmagnet/internal/config"
	"go.uber.org/fx"
)

func NewConfigModule[T any](key string, defaultValue interface{}) fx.Option {
	return fx.Module(
		"config:"+key,
		fx.Provide(
			fx.Annotate(
				func() config.Spec {
					return config.Spec{
						Key:          key,
						DefaultValue: defaultValue,
					}
				},
				fx.ResultTags(`group:"config_specs"`),
			),
			func(r config.ResolvedConfig) (T, error) {
				value, ok := r.NodeMap[key].Value.(T)
				if !ok {
					return value, errors.New("unexpected config type")
				}
				return value, nil
			},
		),
	)
}
