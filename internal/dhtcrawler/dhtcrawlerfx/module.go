package dhtcrawlerfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/dhtcrawler"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		dhtcrawler.Namespace,
		configfx.NewConfigModule[dhtcrawler.Config](dhtcrawler.Namespace, dhtcrawler.NewDefaultConfig()),
		fx.Provide(
			fx.Annotate(
				dhtcrawler.New,
				fx.ResultTags(`group:"worker_options"`, `group:"health_check_options"`),
			),
			dhtcrawler.NewDiscoveredNodes,
		),
	)
}
