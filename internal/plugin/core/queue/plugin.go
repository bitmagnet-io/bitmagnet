package queue

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/queuemetrics"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/handler"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/manager"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/server"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type deps struct {
	fx.In
	Autostart   queue.Autostart
	DaoProvider database.DaoTransactionProvider
	Logger      *zap.Logger
	Handlers    []handler.Handler `group:"queue_handlers"`
}

var (
	Ref = ref.Root.MustSub("queue")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Runs queued background jobs"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithConfig[deps](Ref.MustSub("autostart"), queue.ParamAutostart),
		builder.WithDependencies[deps](
			logging.Ref,
			persister.Ref,
			postgres.Ref,
		),
		builder.WithFxOption[deps](
			fx.Provide(
				manager.New,
				queuemetrics.New,
			),
		),
		builder.WithWorker(
			func(deps deps) (runner.Provider, worker.Option) {
				return server.New(
						deps.DaoProvider,
						deps.Logger.Named(Ref.String()),
						deps.Handlers...,
					),
					worker.Options(
						worker.WithDependencies(
							persister.Ref,
						),
						worker.WithAutostart(bool(deps.Autostart)),
					)
			},
		),
		// builder.WithPrometheusCollector(
		// 	func(cfg config, deps deps) prom.Collector {
		// 		return prometheus.New(deps.DaoProvider, deps.Logger.Named(refMetrics.String()))
		// 	},
		// ),
	)
)
