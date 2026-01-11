//go:build !dev

package provider

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/auth"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/batcher"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dht/crawler"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dht/server"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/health"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server"
	http_server_cors "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server/cors"
	http_server_docs "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server/docs"
	http_server_graphql "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server/graphql"
	http_server_importer "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server/importer"
	http_server_logging "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server/logging"
	http_server_webui "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server/webui"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/info_hash_blocker"
	logging_file_rotator "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging/file_rotator"
	logging_json "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging/json"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/meta_info"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/metrics"
	metrics_prometheus "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/metrics/prometheus"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/persister"
	postgres_cache "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/postgres/cache"
	postgres_migrations "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/postgres/migrator"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/search"
	tmdb "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/tmdb"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/torznab"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/worker"
)

func init() {
	corePlugins = append(corePlugins,
		auth.Plugin,
		batcher.Plugin,
		classifier.Plugin,
		crawler.Plugin,
		dht.Plugin,
		health.Plugin,
		http_server_cors.Plugin,
		http_server_docs.Plugin,
		http_server_graphql.Plugin,
		http_server_importer.Plugin,
		http_server_logging.Plugin,
		http_server_webui.Plugin,
		http_server.Plugin,
		info_hash_blocker.Plugin,
		logging_file_rotator.Plugin,
		logging_json.Plugin,
		meta_info.Plugin,
		metrics_prometheus.Plugin,
		metrics.Plugin,
		persister.Plugin,
		postgres_cache.Plugin,
		postgres_migrations.Plugin,
		processor.Plugin,
		queue.Plugin,
		search.Plugin,
		server.Plugin,
		tmdb.Plugin,
		torznab.Plugin,
		worker.Plugin,
	)
}
