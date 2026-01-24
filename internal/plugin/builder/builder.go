package builder

import (
	"errors"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
	config_registry "github.com/bitmagnet-io/bitmagnet/internal/config/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	workers_registry "github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type builder[Deps any] struct {
	ref           ref.Ref
	activation    plugin.Activation
	dependencies  []ref.Ref
	params        []config_registry.Param
	i18nProviders []i18n.Provider
	errors        ref.Map[error]
	fxOptions     []fx.Option
	commands      []plugin.Command
}

func NewPlugin[Deps any](rf ref.Ref, options ...Option[Deps]) plugin.Plugin {
	b := &builder[Deps]{
		ref:        rf,
		activation: plugin.ActivationAuto,
		errors:     ref.NewMap[error](),
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
		b.errors,
		i18n.Providers(b.i18nProviders...),
		b.commands,
		fx.Options(b.fxOptions...),
	)
}

type Option[Deps any] func(*builder[Deps])

func WithOptions[Deps any](options ...Option[Deps]) Option[Deps] {
	return func(b *builder[Deps]) {
		for _, option := range options {
			option(b)
		}
	}
}

func WithDescription[Deps any](description string) Option[Deps] {
	return func(b *builder[Deps]) {
		WithI18nMessage[Deps](
			b.ref,
			"description for plugin: "+b.ref.String(),
			i18n.WithOther(description),
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

func WithConfig[Deps any, T any](ref ref.Ref, p param.Param[T]) Option[Deps] {
	return WithOptions(
		func(b *builder[Deps]) {
			b.params = append(b.params, config_registry.Param{
				Ref:     ref,
				Plugin:  b.ref,
				Untyped: p,
			})
		},
		WithFxOption[Deps](
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
						return value, fmt.Errorf(
							"failed to cast from %T to %T",
							resolved.Value(),
							value,
						)
					}

					return value, nil
				},
			),
		),
	)
}

func WithI18nMessage[Deps any](ref ref.Ref, description string, options ...i18n.MessageOption) Option[Deps] {
	return WithI18nProvider[Deps](i18n.NewProvider(
		i18n.NewMessage(ref.String(), description, options...),
	))
}

func WithI18nProvider[Deps any](provider i18n.Provider) Option[Deps] {
	return func(builder *builder[Deps]) {
		builder.i18nProviders = append(builder.i18nProviders, provider)
	}
}

func WithError[Deps any](ref ref.Ref, err error) Option[Deps] {
	return func(builder *builder[Deps]) {
		builder.errors.Set(ref, err)
	}
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

func WithAuthObjectActions[Deps any](
	provider func(Deps) []rbac.ObjectAction,
) Option[Deps] {
	return WithFxOption[Deps](
		fx.Provide(
			fx.Annotate(
				func(deps Deps) rbac.ObjectActionProvider {
					return func() []rbac.ObjectAction {
						return provider(deps)
					}
				},
				fx.ResultTags(`group:"auth_object_actions"`),
			),
		),
	)
}

func WithPermissionProvider[Deps any](
	provider func(Deps) []rbac.Permission,
) Option[Deps] {
	return WithFxOption[Deps](
		fx.Provide(
			fx.Annotate(
				func(deps Deps) rbac.PermissionProvider {
					return func() []rbac.Permission {
						return provider(deps)
					}
				},
				fx.ResultTags(`group:"auth_permissions"`),
			),
		),
	)
}
