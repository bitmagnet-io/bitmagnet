package health

import (
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/version"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type (
	config struct{}

	deps struct {
		fx.In
		Checker health.Checker
	}
)

var (
	Ref = core.Ref.MustSub("health")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[config, deps](),
		builder.WithFxOption[config, deps](
			fx.Provide(
				fx.Annotate(
					func(options []health.CheckerOption) health.Checker {
						return health.NewChecker(options...)
					},
					fx.ParamTags(`group:"health_check_options"`),
				),
				fx.Annotate(
					func() health.CheckerOption {
						return health.WithInfo(map[string]any{
							"name":    "bitmagnet",
							"version": version.GitTag,
						})
					},
					fx.ResultTags(`group:"health_check_options"`),
				),
				fx.Annotate(
					health.NewHTTPOption,
					fx.ResultTags(`group:"http_server_options"`),
				),
				fx.Annotate(
					func(checker health.Checker) registry.Option {
						return registry.WithWorker(
							Ref.String(),
							checker.Runner(),
						)
					},
					fx.ResultTags(`group:"worker_options"`),
				),
			),
		),
		builder.WithWorkerRegistryOption(
			func(cfg config, deps deps) registry.Option {
				return registry.WithWorker(Ref.String(), deps.Checker)
			},
		),
		builder.WithGinOption(Ref, func(cfg config, deps deps) gin.OptionFunc {
			return health.NewHTTPOption(deps.Checker)
		}),
	)
)
