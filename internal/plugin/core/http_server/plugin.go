package http_server

import (
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/env"
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver/circuitbreaker"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	circuitbreaker.Handler
	httpserver.LocalAddress
}

var (
	Ref = core.Ref.MustSub("http_server")

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
					func(
						options []httpserver.Option,
					) (gin.OptionFunc, error) {
						return resolveOptions(options)
					},
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
	// paramMap := make(map[string]struct{})
	// for _, p := range optionsNames {
	// 	paramMap[p] = struct{}{}
	// }

	// enabledOptions := make([]httpserver.Option, 0, len(options))

	// foundMap := make(map[string]struct{}, len(options))
	// for _, o := range options {
	// 	if _, ok := foundMap[o.Key()]; ok {
	// 		return nil, fmt.Errorf("duplicate http server option: '%s'", o.Key())
	// 	}

	// 	foundMap[o.Key()] = struct{}{}

	// 	enabled := false
	// 	if _, ok := paramMap["*"]; ok {
	// 		enabled = true
	// 	} else if _, ok := paramMap[o.Key()]; ok {
	// 		enabled = true
	// 	}

	// 	if enabled {
	// 		enabledOptions = append(enabledOptions, o)
	// 	}
	// }

	// for p := range paramMap {
	// 	if _, ok := foundMap[p]; !ok && p != "*" {
	// 		return nil, fmt.Errorf("unknown http server option: '%s'", p)
	// 	}
	// }

	slices.SortFunc(options, func(a, b httpserver.Option) int {
		return a.Compare(b)
	})

	return func(engine *gin.Engine) {
		for _, opt := range options {
			opt.Apply(engine)
		}
	}, nil
}
