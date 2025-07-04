package queuefx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/blocking"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/manager"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/prometheus"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/server"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	libprometheus "github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New() fx.Option {
	return fx.Module(
		"queue",
		fx.Provide(
			fx.Annotate(
				func(
					handlers []lazy.Lazy[handler.Handler],
					daoProvider database.DaoTransactionProvider,
					logger *zap.SugaredLogger,
				) registry.Option {
					return registry.WithWorker(
						"queue_server",
						server.New(
							daoProvider,
							logger.Named("queue_server"),
							handlers...,
						),
						worker.WithDependencies(database.Namespace, blocking.Namespace),
					)
				},
				fx.ParamTags(`group:"queue_handlers"`),
				fx.ResultTags(`group:"worker_options"`),
			),
			fx.Annotate(
				func(
					daoProvider database.DaoProvider,
					logger *zap.SugaredLogger,
				) libprometheus.Collector {
					return prometheus.New(
						daoProvider,
						logger.Named("queue_metrics_collector"),
					)
				},
				fx.ResultTags(`group:"prometheus_collectors"`),
			),
			manager.New,
		),
	)
}
