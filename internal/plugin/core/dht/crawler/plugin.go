package crawler

import (
	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/dhtcrawler"
	"github.com/bitmagnet-io/bitmagnet/internal/dhtcrawler/dhtcrawlerhealthcheck"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/info_hash_blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dht"
	plugin_dht_server "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dht/server"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/meta_info"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline/processor"
	plugin_worker "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/worker"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	Autostart      dhtcrawler.Autostart
	Runner         dhtcrawler.Runner
	LastResponses  *atomic.Value[server.LastResponses]
	WorkerRegistry registry.StateProvider
}

var (
	Ref = dht.Ref.MustSub("crawler")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Crawls the DHT for torrents"),
		builder.WithActivation[deps](plugin.ActivationEnabled),
		builder.WithDependencies[deps](
			classifier.Ref,
			info_hash_blocker.Ref,
			logging.Ref,
			meta_info.Ref,
			persister.Ref,
			plugin_dht_server.Ref,
			postgres.Ref,
			processor.Ref,
			plugin_worker.Ref,
		),
		builder.WithConfig[deps](Ref.MustSub("autostart"), dhtcrawler.ParamAutostart),
		builder.WithConfig[deps](Ref.MustSub("bootstrap_nodes"), dhtcrawler.ParamBootstrapNodes),
		builder.WithConfig[deps](Ref.MustSub("reseed_bootstrap_nodes_interval"), dhtcrawler.ParamReseedBootstrapNodesInterval),
		builder.WithConfig[deps](Ref.MustSub("save_files_threshold"), dhtcrawler.ParamSaveFilesThreshold),
		builder.WithConfig[deps](Ref.MustSub("save_pieces"), dhtcrawler.ParamSavePieces),
		builder.WithConfig[deps](Ref.MustSub("rescrape_threshold"), dhtcrawler.ParamRescrapeThreshold),
		builder.WithFxOption[deps](
			fx.Provide(
				fx.Private,
				func(metrics *metrics.Registry) (*metrics.Component, error) {
					return metrics.NewComponent(Ref)
				},
			),
			fx.Provide(
				dhtcrawler.New,
			),
		),
		builder.WithWorker(
			func(deps deps) (runner.Provider, worker.Option) {
				return deps.Runner, worker.Options(
					worker.WithDependencies(
						info_hash_blocker.Ref,
						postgres.Ref,
						plugin_dht_server.Ref,
					),
					worker.WithAutostart(bool(deps.Autostart)),
				)
			},
		),
		builder.WithHealthCheckerOption(func(deps deps) health.CheckerOption {
			return dhtcrawlerhealthcheck.New(Ref.String(), func() bool {
				return deps.WorkerRegistry.WorkersState().Get(Ref).State != worker.StateIdle
			}, deps.LastResponses)
		}),
	)
)
