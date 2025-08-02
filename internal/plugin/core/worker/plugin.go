package worker

import (
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/workershealthcheck"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	config struct{}

	deps struct {
		fx.In
		HealthCheck *workershealthcheck.WorkersHealthCheck
	}
)

var (
	Ref = core.Ref.MustSub("worker")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[config, deps](),
		builder.WithFxOption[config, deps](
			fx.Provide(
				fx.Annotate(
					func(
						options []registry.Option,
						logger *zap.Logger,
					) (*registry.Registry, error) {
						return registry.NewRegistry(
							logger.Named(Ref.String()),
							options...,
						)
					},
					fx.ParamTags(`group:"worker_options"`),
				),
				workershealthcheck.New,
				fx.Annotate(
					registry.NewCircuitBreaker,
					fx.As(new(registry.StateProvider)),
					fx.As(fx.Self()),
				),
			),
			fx.Invoke(
				func(r *registry.Registry, cb registry.CircuitBreaker) error {
					return cb.ReceiveRegistry(r)
				},
			),
		),
		builder.WithHealthCheckerOption(
			func(_ config, deps deps) health.CheckerOption {
				return health.WithCheck(deps.HealthCheck.Check())
			},
		),
		builder.WithCliCommand[config, deps](NewStartCommand()),
	)
)
