package appfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/app/cmd/searchcmd"
	"github.com/bitmagnet-io/bitmagnet/internal/app/cmd/torrentcmd"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/app/boilerplateappfx"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver/httpserverfx"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classifierfx"
	"github.com/bitmagnet-io/bitmagnet/internal/database/databasefx"
	"github.com/bitmagnet-io/bitmagnet/internal/database/migrations"
	"github.com/bitmagnet-io/bitmagnet/internal/dhtcrawler/dhtcrawlerfx"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlfx"
	"github.com/bitmagnet-io/bitmagnet/internal/importer/importerfx"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/dhtfx"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainfofx"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/queuefx"
	"github.com/bitmagnet-io/bitmagnet/internal/redis/redisfx"
	"github.com/bitmagnet-io/bitmagnet/internal/telemetry/telemetryfx"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab/torznabfx"
	"github.com/bitmagnet-io/bitmagnet/internal/version/versionfx"
	"github.com/bitmagnet-io/bitmagnet/internal/webui"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"app",
		boilerplateappfx.New(),
		classifierfx.New(),
		dhtcrawlerfx.New(),
		dhtfx.New(),
		databasefx.New(),
		gqlfx.New(),
		httpserverfx.New(),
		importerfx.New(),
		metainfofx.New(),
		queuefx.New(),
		redisfx.New(),
		telemetryfx.New(),
		torznabfx.New(),
		versionfx.New(),
		// cli commands:
		fx.Provide(
			searchcmd.New,
			torrentcmd.New,
		),
		fx.Provide(webui.New),
		fx.Decorate(migrations.NewDecorator),
	)
}
