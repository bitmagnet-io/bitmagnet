package http_server

import (
	"fmt"
	"sort"

	"github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver/circuitbreaker"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type (
	Config = httpserver.Config

	deps struct {
		fx.In
		circuitbreaker.Handler
	}
)

var (
	Ref = core.Ref.MustSub("http_server")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[Config, deps](),
		builder.WithDependencies[Config, deps](
			config.Ref,
		),
		builder.WithDefaultConfig[Config, deps](httpserver.NewDefaultConfig()),
		builder.WithFxOption[Config, deps](
			fx.Provide(
				circuitbreaker.New,
				fx.Annotate(
					func(
						options []httpserver.Option,
						config httpserver.Config,
					) (gin.OptionFunc, error) {
						return resolveOptions(options, config.Options)
					},
					fx.ParamTags(`group:"http_server_options"`),
				),
			),
			fx.Invoke(
				func(
					config httpserver.Config,
					option gin.OptionFunc,
					handler circuitbreaker.Handler,
				) error {
					gin.SetMode(config.GinMode)

					return handler.SetOption(option)
				},
			),
		),
		builder.WithWorkerRegistryOption(func(cfg Config, deps deps) registry.Option {
			return registry.WithWorker(
				Ref.String(),
				httpserver.New(deps, cfg.LocalAddress),
				worker.WithAutostart(),
			)
		}),
	)
)

func resolveOptions(options []httpserver.Option, optionsNames []string) (gin.OptionFunc, error) {
	paramMap := make(map[string]struct{})
	for _, p := range optionsNames {
		paramMap[p] = struct{}{}
	}

	enabledOptions := make([]httpserver.Option, 0, len(options))

	foundMap := make(map[string]struct{}, len(options))
	for _, o := range options {
		if _, ok := foundMap[o.Key()]; ok {
			return nil, fmt.Errorf("duplicate http server option: '%s'", o.Key())
		}

		foundMap[o.Key()] = struct{}{}

		enabled := false
		if _, ok := paramMap["*"]; ok {
			enabled = true
		} else if _, ok := paramMap[o.Key()]; ok {
			enabled = true
		}

		if enabled {
			enabledOptions = append(enabledOptions, o)
		}
	}

	for p := range paramMap {
		if _, ok := foundMap[p]; !ok && p != "*" {
			return nil, fmt.Errorf("unknown http server option: '%s'", p)
		}
	}

	sort.Slice(enabledOptions, func(i, j int) bool {
		return enabledOptions[i].Key() < enabledOptions[j].Key()
	})

	return func(engine *gin.Engine) {
		for _, opt := range enabledOptions {
			opt.Apply(engine)
		}
	}, nil
}
