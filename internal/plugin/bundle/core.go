package bundle

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	database_gorm_cache "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/gorm/cache"
	database_info_hash_blocker "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/info_hash_blocker"
	database_migrations "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/migrator"
	database_postgres "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	database_search "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/search"
	dht_crawler "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dht/crawler"
	dht_server "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dht/server"
	dht_socket "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dht/socket"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/health"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server"
	http_server_cors "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server/cors"
	http_server_docs "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server/docs"
	http_server_graphql "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server/graphql"
	http_server_importer "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server/importer"
	http_server_logging "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server/logging"
	http_server_webui "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server/webui"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	logging_console "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging/console"
	logging_file_rotator "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging/file_rotator"
	logging_json "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging/json"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/meta_info"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/metrics"
	metrics_prometheus "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/metrics/prometheus"
	pipeline_batcher "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline/batcher"
	pipeline_classifier "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline/classifier"
	pipeline_persister "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline/persister"
	pipeline_indexer "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/queue"
	tmdb_compat "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/tmdb/compat"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/torznab"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/worker"
)

var Core = MustNew(
	core.Ref,
	config.Plugin,
	database_gorm_cache.Plugin,
	database_info_hash_blocker.Plugin,
	database_migrations.Plugin,
	database_postgres.Plugin,
	database_search.Plugin,
	dht_crawler.Plugin,
	dht_server.Plugin,
	dht_socket.Plugin,
	health.Plugin,
	http_server.Plugin,
	http_server_cors.Plugin,
	http_server_docs.Plugin,
	http_server_graphql.Plugin,
	http_server_importer.Plugin,
	http_server_logging.Plugin,
	http_server_webui.Plugin,
	i18n.Plugin,
	logging.Plugin,
	logging_console.Plugin,
	logging_file_rotator.Plugin,
	logging_json.Plugin,
	meta_info.Plugin,
	metrics.Plugin,
	metrics_prometheus.Plugin,
	pipeline_batcher.Plugin,
	pipeline_classifier.Plugin,
	pipeline_persister.Plugin,
	pipeline_indexer.Plugin,
	queue.Plugin,
	// tmdb.Plugin,
	tmdb_compat.Plugin,
	torznab.Plugin,
	worker.Plugin,
)
