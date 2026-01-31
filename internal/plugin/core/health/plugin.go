package health

import (
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/version"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	Checker health.Checker
}

var (
	Ref = ref.Root.MustSub("health")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides health check functionality"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithFxOption[deps](
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
				// fx.Annotate(
				// 	func(checker health.Checker) registry.Option {
				// 		return registry.WithWorker(
				// 			Ref,
				// 			checker.Runner(),
				// 		)
				// 	},
				// 	fx.ResultTags(`group:"worker_options"`),
				// ),
			),
		),
		builder.WithWorker(
			func(deps deps) (runner.Provider, worker.Option) {
				return deps.Checker, nil
			},
		),
		builder.WithGinOption(Ref, 0, func(deps deps) gin.OptionFunc {
			return health.NewHTTPOption(deps.Checker)
		}),
	)
)
