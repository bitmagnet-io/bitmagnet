package httpserverfx

import (
	"fmt"
	"sort"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/blocking"
	"github.com/bitmagnet-io/bitmagnet/internal/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver/circuitbreaker"
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver/cors"
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver/ginzap"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New() fx.Option {
	return fx.Module(
		httpserver.Namespace,
		configfx.NewConfigModule[httpserver.Config](httpserver.Namespace, httpserver.NewDefaultConfig()),
		fx.Provide(
			circuitbreaker.New,
			fx.Annotate(
				func(logger *zap.Logger) httpserver.Option {
					return httpserver.NewOption("logger", func(engine *gin.Engine) {
						engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
					})
				},
				fx.ResultTags(`group:"http_server_options"`),
			),
			fx.Annotate(
				func() httpserver.Option {
					return httpserver.NewOption("recovery", func(engine *gin.Engine) {
						engine.Use(gin.Recovery())
					})
				},
				fx.ResultTags(`group:"http_server_options"`),
			),
			fx.Annotate(
				func(
					options []httpserver.Option,
					config httpserver.Config,
				) (gin.OptionFunc, error) {
					return resolveOptions(options, config.Options)
				},
				fx.ParamTags(`group:"http_server_options"`),
			),
			fx.Annotate(
				func(
					config httpserver.Config,
					handler circuitbreaker.Handler,
				) (registry.Option, error) {
					return registry.WithWorker(
						httpserver.Namespace,
						httpserver.New(handler, config.LocalAddress),
						worker.WithDependencies(
							database.Namespace,
							blocking.Namespace,
							health.Namespace,
						),
					), nil
				},
				fx.ResultTags(`group:"worker_options"`),
			),
			fx.Annotate(
				func(
					config httpserver.Config,
					logger *zap.SugaredLogger,
				) httpserver.Option {
					return cors.New(config.Cors, logger)
				},
				fx.ResultTags(`group:"http_server_options"`),
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
	)
}

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
