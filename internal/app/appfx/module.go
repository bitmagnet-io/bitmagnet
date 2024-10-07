package appfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/app/cmd/classifiercmd"
	"github.com/bitmagnet-io/bitmagnet/internal/app/cmd/processcmd"
	"github.com/bitmagnet-io/bitmagnet/internal/app/cmd/reprocesscmd"
	"github.com/bitmagnet-io/bitmagnet/internal/blocking/blockingfx"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/app/boilerplateappfx"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver/httpserverfx"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classifierfx"
	"github.com/bitmagnet-io/bitmagnet/internal/database/databasefx"
	"github.com/bitmagnet-io/bitmagnet/internal/database/migrations"
	"github.com/bitmagnet-io/bitmagnet/internal/dhtcrawler/dhtcrawlerfx"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlfx"
	"github.com/bitmagnet-io/bitmagnet/internal/health/healthfx"
	"github.com/bitmagnet-io/bitmagnet/internal/importer/importerfx"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/metricsfx"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/processorfx"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/dhtfx"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainfofx"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/queuefx"
	"github.com/bitmagnet-io/bitmagnet/internal/telemetry/telemetryfx"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb/tmdbfx"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab/torznabfx"
	"github.com/bitmagnet-io/bitmagnet/internal/version/versionfx"
	"github.com/bitmagnet-io/bitmagnet/internal/webui"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"app",
		blockingfx.New(),
		boilerplateappfx.New(),
		dhtcrawlerfx.New(),
		dhtfx.New(),
		databasefx.New(),
		gqlfx.New(),
		healthfx.New(),
		httpserverfx.New(),
		importerfx.New(),
		metainfofx.New(),
		metricsfx.New(),
		processorfx.New(),
		queuefx.New(),
		telemetryfx.New(),
		tmdbfx.New(),
		torznabfx.New(),
		versionfx.New(),
		classifierfx.New(),
		// cli commands:
		fx.Provide(
			classifiercmd.New,
			reprocesscmd.New,
			processcmd.New,
		),
		fx.Provide(webui.New),
		fx.Decorate(migrations.NewDecorator),
	)
}
