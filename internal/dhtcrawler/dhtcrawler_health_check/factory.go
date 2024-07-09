package dhtcrawler_health_check

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/dhtcrawler"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"go.uber.org/fx"
	"time"
)

type Params struct {
	fx.In
	Config           dhtcrawler.Config
	DhtCrawlerActive *concurrency.AtomicValue[bool] `name:"dht_crawler_active"`
	Client           lazy.Lazy[client.Client]
}

type Result struct {
	fx.Out
	Option health.CheckerOption `group:"health_check_options"`
}

func New(params Params) Result {
	return Result{
		Option: health.WithPeriodicCheck(
			time.Minute*10,
			time.Second*10,
			NewCheck(params.DhtCrawlerActive, params.Client, params.Config.BootstrapNodes),
		),
	}
}
