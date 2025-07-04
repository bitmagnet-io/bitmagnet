package healthfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/version"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"health",
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
						health.Namespace,
						checker.Runner,
					)
				},
				fx.ResultTags(`group:"worker_options"`),
			),
		),
	)
}
