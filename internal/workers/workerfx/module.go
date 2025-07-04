package workerfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/workers"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/workershealthcheck"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New() fx.Option {
	return fx.Module(
		workers.Namespace,
		fx.Provide(
			fx.Annotate(
				func(
					options []registry.Option,
					logger *zap.SugaredLogger,
				) (*registry.Registry, error) {
					return registry.NewRegistry(
						logger.Named(worker.Namespace),
						options...,
					)
				},
				fx.ParamTags(`group:"worker_options"`),
			),
			workershealthcheck.New,
			fx.Annotate(
				func(hc *workershealthcheck.WorkersHealthCheck) health.CheckerOption {
					return health.WithCheck(hc.Check())
				},
				fx.ResultTags(`group:"health_check_options"`),
			),
		),
		fx.Invoke(
			func(r *registry.Registry, hc *workershealthcheck.WorkersHealthCheck) error {
				return hc.SetRegistry(r)
			},
		),
	)
}
