package config

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config"
	"github.com/bitmagnet-io/bitmagnet/internal/config/lookup"
	"github.com/bitmagnet-io/bitmagnet/internal/config/manager"
	"github.com/bitmagnet-io/bitmagnet/internal/config/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/fs"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/validation"
	"github.com/spf13/afero"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
}

var (
	Ref = core.Ref.MustSub("config")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[deps](),
		builder.WithDependencies[deps](
			validation.Ref,
		),
		builder.WithFxOption[deps](
			fx.Provide(
				lookup.NewFromEnv,
				resolver.New,
				func(provider fs.FSProvider) manager.FS {
					return afero.Afero{
						Fs: afero.NewBasePathFs(provider.FSData(), config.SubpathPersisted),
					}
				},
				fx.Private,
			),
			fx.Provide(
				fx.Annotate(
					registry.New,
					fx.ParamTags(`group:"config_registry_options"`),
				),
				func(resolver resolver.Resolver) (resolver.Resolved, error) {
					return resolver.Resolve()
				},
				manager.New,
			),
		),
		// builder.WithCliCommand[deps](
		// 	NewConfigCommand(),
		// ),
		// builder.WithFxOption[config, deps](configfx.New()),
	)
)
