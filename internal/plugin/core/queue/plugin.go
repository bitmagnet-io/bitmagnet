package queue

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/queuemetrics"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/manager"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/prometheus"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/server"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	prom "github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	config struct{}

	deps struct {
		fx.In
		DaoProvider database.DaoTransactionProvider
		Logger      *zap.SugaredLogger
		Handlers    []handler.Handler `group:"queue_handlers"`
	}
)

var (
	Ref = core.Ref.MustSub("queue")

	refMetrics = Ref.MustSub("metrics")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[config, deps](),
		builder.WithDependencies[config, deps](
			logging.Ref,
			persister.Ref,
			postgres.Ref,
		),
		builder.WithFxOption[config, deps](
			fx.Provide(
				manager.New,
				queuemetrics.New,
			),
		),
		builder.WithWorkerRegistryOption(
			func(cfg config, deps deps) registry.Option {
				return registry.WithWorker(
					Ref.String(),
					server.New(
						deps.DaoProvider,
						deps.Logger.Named(Ref.String()),
						deps.Handlers...,
					),
					worker.WithDependencies(
						persister.Ref.String(),
					),
					worker.WithAutostart(),
				)
			},
		),
		builder.WithPrometheusCollector(
			func(cfg config, deps deps) prom.Collector {
				return prometheus.New(deps.DaoProvider, deps.Logger.Named(refMetrics.String()))
			},
		),
	)
)
