package builder

import (
	"errors"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
	config_registry "github.com/bitmagnet-io/bitmagnet/internal/config/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	workers_registry "github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type builder[Deps any] struct {
	ref              ref.Ref
	enabledByDefault bool
	dependencies     []ref.Ref
	options          []fx.Option
	commands         []plugin.Command
}

func CreatePlugin[Deps any](ref ref.Ref, options ...Option[Deps]) plugin.Plugin {
	b := &builder[Deps]{
		ref: ref,
	}

	for _, option := range options {
		option(b)
	}

	// WithFxOption[Deps](configfx.NewConfigModule[Config](ref, b.config))(b)

	return plugin.NewPlugin(ref, b.enabledByDefault, b.dependencies, fx.Options(b.options...), b.commands)
}

type Option[Deps any] func(*builder[Deps])

func WithDependencies[Deps any](dependencies ...ref.Ref) Option[Deps] {
	return func(b *builder[Deps]) {
		b.dependencies = append(b.dependencies, dependencies...)
	}
}

func WithEnabledByDefault[Deps any]() Option[Deps] {
	return func(b *builder[Deps]) {
		b.enabledByDefault = true
	}
}

func WithConfigParam[Deps any, T any](ref ref.Ref, param param.Param[T]) Option[Deps] {
	return WithFxOption[Deps](
		fx.Supply(
			fx.Annotate(
				config_registry.WithParam(ref, param),
				fx.ResultTags(`group:"config_registry_options"`),
			),
		),
		fx.Provide(
			func(allResolved resolver.Resolved) (T, error) {
				var (
					resolved *resolver.Param
					value    T
					ok       bool
				)

				resolved, ok = allResolved.Param(ref)
				if !ok {
					return value, errors.New("missing from resolved")
				}

				value, ok = resolved.Value().(T)
				if !ok {
					return value, errors.New("cast failed")
				}

				return value, nil
			},
		),
	)
}

func WithFxOption[Deps any](options ...fx.Option) Option[Deps] {
	return func(b *builder[Deps]) {
		b.options = append(b.options, options...)
	}
}

func WithGinOption[Deps any](ref ref.Ref, provider func(Deps) gin.OptionFunc) Option[Deps] {
	return WithFxOption[Deps](fx.Provide(
		fx.Annotate(
			func(deps Deps) httpserver.Option {
				return httpserver.NewOption(ref.String(), provider(deps))
			},
			fx.ResultTags(`group:"http_server_options"`),
		),
	))
}

func WithGormPlugin[Deps any](provider func(Deps) gorm.Plugin) Option[Deps] {
	return WithFxOption[Deps](fx.Provide(
		fx.Annotate(
			provider,
			fx.ResultTags(`group:"gorm_plugins"`),
		),
	))
}

func WithZapCore[Deps any](provider func(Deps) zapcore.Core) Option[Deps] {
	return WithFxOption[Deps](fx.Provide(
		fx.Annotate(
			provider,
			fx.ResultTags(`group:"zap_cores"`),
		),
	))
}

func WithWorkerRegistryOption[Deps any](provider func(Deps) workers_registry.Option) Option[Deps] {
	return WithFxOption[Deps](fx.Provide(
		fx.Annotate(
			provider,
			fx.ResultTags(`group:"worker_options"`),
		),
	))
}

func WithHealthCheckerOption[Deps any](provider func(Deps) health.CheckerOption) Option[Deps] {
	return WithFxOption[Deps](fx.Provide(
		fx.Annotate(
			provider,
			fx.ResultTags(`group:"health_check_options"`),
		),
	))
}

func WithQueueHandler[Deps any](provider func(Deps) handler.Handler) Option[Deps] {
	return WithFxOption[Deps](fx.Provide(
		fx.Annotate(
			provider,
			fx.ResultTags(`group:"queue_handlers"`),
		),
	))
}

func WithPrometheusCollector[Deps any](provider func(Deps) prometheus.Collector) Option[Deps] {
	return WithFxOption[Deps](fx.Provide(
		fx.Annotate(
			provider,
			fx.ResultTags(`group:"prometheus_collectors"`),
		),
	))
}

func WithCliCommand[Deps any](
	commands ...plugin.Command,
) Option[Deps] {
	return func(b *builder[Deps]) {
		b.commands = append(b.commands, commands...)
	}
}
