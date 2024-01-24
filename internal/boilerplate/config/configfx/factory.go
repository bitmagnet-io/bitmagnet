package configfx

import (
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config"
	"go.uber.org/fx"
)

// NewConfigModule
//
// Linter false positive: `Error: cannot infer T (internal/boilerplate/config/configfx/factory.go:9:22) (typecheck)`
//
//nolint:typecheck
func NewConfigModule[T any](key string, defaultValue interface{}) fx.Option {
	return fx.Module(
		"config:"+key,
		fx.Provide(
			fx.Annotated{
				Group: "config_specs",
				Target: func() config.Spec {
					return config.Spec{
						Key:          key,
						DefaultValue: defaultValue,
					}
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Target: func(r config.ResolvedConfig) (cfg T, err error) {
					v, ok := r.NodeMap[key].Value.(T)
					if !ok {
						err = errors.New("unexpected config type")
						return
					}
					return v, nil
				},
			},
		),
	)
}
