package builder

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/configfx"
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

type builder[Config any, Deps any] struct {
	ref              ref.Ref
	enabledByDefault bool
	dependencies     []ref.Ref
	config           Config
	options          []fx.Option
	commands         []plugin.Command
}

func CreatePlugin[Config any, Deps any](ref ref.Ref, options ...Option[Config, Deps]) plugin.Plugin {
	b := &builder[Config, Deps]{
		ref: ref,
	}

	for _, option := range options {
		option(b)
	}

	WithFxOption[Config, Deps](configfx.NewConfigModule[Config](ref.String(), b.config))(b)

	return plugin.NewPlugin(ref, b.enabledByDefault, b.dependencies, fx.Options(b.options...), b.commands)
}

type Option[Config any, Deps any] func(*builder[Config, Deps])

func WithDependencies[Config any, Deps any](dependencies ...ref.Ref) Option[Config, Deps] {
	return func(b *builder[Config, Deps]) {
		b.dependencies = append(b.dependencies, dependencies...)
	}
}

func WithEnabledByDefault[Config any, Deps any]() Option[Config, Deps] {
	return func(b *builder[Config, Deps]) {
		b.enabledByDefault = true
	}
}

func WithDefaultConfig[Config any, Deps any](config Config) Option[Config, Deps] {
	return func(b *builder[Config, Deps]) {
		b.config = config
	}
}

func WithFxOption[Config any, Deps any](options ...fx.Option) Option[Config, Deps] {
	return func(b *builder[Config, Deps]) {
		b.options = append(b.options, options...)
	}
}

func WithGinOption[Config any, Deps any](ref ref.Ref, provider func(Config, Deps) gin.OptionFunc) Option[Config, Deps] {
	return WithFxOption[Config, Deps](fx.Provide(
		fx.Annotate(
			func(cfg Config, deps Deps) httpserver.Option {
				return httpserver.NewOption(ref.String(), provider(cfg, deps))
			},
			fx.ResultTags(`group:"http_server_options"`),
		),
	))
}

func WithGormPlugin[Config any, Deps any](provider func(Config, Deps) gorm.Plugin) Option[Config, Deps] {
	return WithFxOption[Config, Deps](fx.Provide(
		fx.Annotate(
			provider,
			fx.ResultTags(`group:"gorm_plugins"`),
		),
	))
}

func WithZapCore[Config any, Deps any](provider func(Config, Deps) zapcore.Core) Option[Config, Deps] {
	return WithFxOption[Config, Deps](fx.Provide(
		fx.Annotate(
			provider,
			fx.ResultTags(`group:"zap_cores"`),
		),
	))
}

func WithWorkerRegistryOption[Config any, Deps any](provider func(Config, Deps) workers_registry.Option) Option[Config, Deps] {
	return WithFxOption[Config, Deps](fx.Provide(
		fx.Annotate(
			provider,
			fx.ResultTags(`group:"worker_options"`),
		),
	))
}

func WithHealthCheckerOption[Config any, Deps any](provider func(Config, Deps) health.CheckerOption) Option[Config, Deps] {
	return WithFxOption[Config, Deps](fx.Provide(
		fx.Annotate(
			provider,
			fx.ResultTags(`group:"health_check_options"`),
		),
	))
}

func WithQueueHandler[Config any, Deps any](provider func(Config, Deps) handler.Handler) Option[Config, Deps] {
	return WithFxOption[Config, Deps](fx.Provide(
		fx.Annotate(
			provider,
			fx.ResultTags(`group:"queue_handlers"`),
		),
	))
}

func WithPrometheusCollector[Config any, Deps any](provider func(Config, Deps) prometheus.Collector) Option[Config, Deps] {
	return WithFxOption[Config, Deps](fx.Provide(
		fx.Annotate(
			provider,
			fx.ResultTags(`group:"prometheus_collectors"`),
		),
	))
}

func WithCliCommand[Config any, Deps any](
	commands ...plugin.Command,
) Option[Config, Deps] {
	return func(b *builder[Config, Deps]) {
		b.commands = append(b.commands, commands...)
	}
}
