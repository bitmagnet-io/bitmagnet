package errors

import (
	"github.com/bitmagnet-io/bitmagnet/internal/error_registry"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	Registry error_registry.Registry
}

var (
	Ref = core.Ref.MustSub("errors")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithFxOption[deps](
			fx.Provide(
				fx.Annotate(
					func(options []error_registry.Option) error_registry.Registry {
						return error_registry.New(options...)
					},
					fx.ParamTags(`group:"error_registry_options"`),
				),
			),
		),
	)
)
