package worker

import (
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/workershealthcheck"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type deps struct {
	fx.In
	HealthCheck *workershealthcheck.WorkersHealthCheck
}

var (
	Ref = ref.Root.MustSub("worker")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Manages workers"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithFxOption[deps](
			fx.Provide(
				fx.Annotate(
					func(
						options []registry.Option,
						logger *zap.Logger,
					) (*registry.Registry, error) {
						return registry.NewRegistry(
							logger,
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
			func(deps deps) health.CheckerOption {
				return health.WithCheck(deps.HealthCheck.Check())
			},
		),
		builder.WithCliCommand[deps](NewStartCommand()),
	)
)
