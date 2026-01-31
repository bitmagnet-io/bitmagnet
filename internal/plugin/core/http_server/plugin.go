package http_server

import (
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver/circuitbreaker"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	circuitbreaker.Handler
	httpserver.LocalAddress
}

var (
	Ref = ref.Root.MustSub("http_server")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Runs the HTTP server"),
		builder.WithActivation[deps](plugin.ActivationEnabled),
		builder.WithDependencies[deps](
			config.Ref,
		),
		builder.WithConfig[deps](Ref.MustSub("local_address"), httpserver.ParamLocalAddress),
		builder.WithGinOption(
			Ref.MustSub("recovery"),
			httpserver.PhasePre,
			func(deps) gin.OptionFunc {
				return func(engine *gin.Engine) {
					engine.Use(gin.Recovery())
				}
			},
		),
		builder.WithGinOption(
			Ref.MustSub("context"),
			httpserver.PhasePre,
			func(deps) gin.OptionFunc {
				return func(e *gin.Engine) {
					e.Use(httpserver.GinContextToContextMiddleware())
				}
			},
		),
		builder.WithFxOption[deps](
			fx.Provide(
				circuitbreaker.New,
				fx.Annotate(
					resolveOptions,
					fx.ParamTags(`group:"http_server_options"`),
				),
			),
			fx.Invoke(
				func(
					env env.VarsLookup,
					option gin.OptionFunc,
					handler circuitbreaker.Handler,
				) error {
					// Set gin mode if it wasn't already
					if _, ok := env.LookupVar(gin.EnvGinMode); !ok {
						gin.SetMode(gin.ReleaseMode)
					}

					return handler.SetOption(option)
				},
			),
		),
		builder.WithWorker(
			func(deps deps) (runner.Provider, worker.Option) {
				return httpserver.New(deps.Handler, deps.LocalAddress),
					worker.WithAutostart(true)
			},
		),
	)
)

func resolveOptions(options []httpserver.Option) (gin.OptionFunc, error) {
	slices.SortFunc(options, func(a, b httpserver.Option) int {
		return a.Compare(b)
	})

	return func(engine *gin.Engine) {
		for _, opt := range options {
			opt.Apply(engine)
		}
	}, nil
}
