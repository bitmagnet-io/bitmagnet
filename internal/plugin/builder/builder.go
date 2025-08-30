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
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type builder[Deps any] struct {
	ref          ref.Ref
	description  string
	activation   plugin.Activation
	dependencies []ref.Ref
	params       []config_registry.Param
	fxOptions    []fx.Option
	commands     []plugin.Command
}

// var paramEnabled = param.MustNew(
// 	param.Description[bool]("Enabled"),
// 	param.Bool[bool](),
// 	param.Default(plugin.EnabledByDefault()),
// )

func NewPlugin[Deps any](rf ref.Ref, options ...Option[Deps]) plugin.Plugin {
	b := &builder[Deps]{
		ref:        rf,
		activation: plugin.ActivationAuto,
	}

	for _, option := range options {
		option(b)
	}

	var activationRef ref.Nullable

	params := b.params

	if b.activation != plugin.ActivationAlways {
		activationRef = ref.NewNullable(rf.MustSub(plugin.KeyActivation))
		params = append(params, config_registry.Param{
			Ref:    activationRef.Ref,
			Plugin: rf,
			Untyped: param.MustNew(
				param.Description[plugin.Activation]("Activation"),
				param.EnumValues(
					plugin.ActivationEnabled,
					plugin.ActivationDisabled,
					plugin.ActivationAuto,
				),
				param.Default(b.activation),
			),
		})
	}

	return plugin.NewPlugin(
		rf,
		activationRef,
		b.dependencies,
		params,
		b.commands,
		fx.Options(b.fxOptions...),
	)
}

type Option[Deps any] func(*builder[Deps])

func options[Deps any](options ...Option[Deps]) Option[Deps] {
	return func(b *builder[Deps]) {
		for _, option := range options {
			option(b)
		}
	}
}

func WithDescription[Deps any](description string) Option[Deps] {
	return func(b *builder[Deps]) {
		b.description = description

		WithFxOption[Deps](
			fx.Supply(
				fx.Annotate(
					&i18n.Message{
						ID:          b.ref.String(),
						Description: "description for plugin " + b.ref.String(),
						Other:       description,
					},
					fx.ResultTags(`group:"i18n_messages"`),
				),
			),
		)(b)
	}
}

func WithDependencies[Deps any](dependencies ...ref.Ref) Option[Deps] {
	return func(b *builder[Deps]) {
		b.dependencies = append(b.dependencies, dependencies...)
	}
}

func WithActivation[Deps any](activation plugin.Activation) Option[Deps] {
	return func(b *builder[Deps]) {
		b.activation = activation
	}
}

func WithConfig[Deps any, T any](ref ref.Ref, param param.Param[T]) Option[Deps] {
	return options(
		func(b *builder[Deps]) {
			b.params = append(b.params, config_registry.Param{
				Ref:     ref,
				Plugin:  b.ref,
				Untyped: param,
			})
		},
		WithFxOption[Deps](
			fx.Supply(
				fx.Annotate(
					&i18n.Message{
						ID:          ref.String(),
						Description: "description for config " + ref.String(),
						Other:       param.Description(),
					},
					fx.ResultTags(`group:"i18n_messages"`),
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
		),
	)
}

func WithFxOption[Deps any](options ...fx.Option) Option[Deps] {
	return func(b *builder[Deps]) {
		b.fxOptions = append(b.fxOptions, options...)
	}
}

func WithGinOption[Deps any](ref ref.Ref, phase httpserver.Phase, provider func(Deps) gin.OptionFunc) Option[Deps] {
	return WithFxOption[Deps](fx.Provide(
		fx.Annotate(
			func(deps Deps) httpserver.Option {
				return httpserver.NewOption(ref.String(), phase, provider(deps))
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

func WithWorker[Deps any](provider func(Deps) (runner.Provider, worker.Option)) Option[Deps] {
	return func(b *builder[Deps]) {
		b.fxOptions = append(b.fxOptions, fx.Provide(
			fx.Annotate(
				func(deps Deps) workers_registry.Option {
					runner, option := provider(deps)

					if option == nil {
						option = worker.Options()
					}

					return workers_registry.WithWorker(
						b.ref,
						runner,
						option,
					)
				},
				fx.ResultTags(`group:"worker_options"`),
			),
		))
	}
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
