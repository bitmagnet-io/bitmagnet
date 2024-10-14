package dhtcrawler_health_check

import (
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
	"go.uber.org/fx"
	"time"
)

type Params struct {
	fx.In
	DhtCrawlerActive       *concurrency.AtomicValue[bool]                 `name:"dht_crawler_active"`
	DhtServerLastResponses *concurrency.AtomicValue[server.LastResponses] `name:"dht_server_last_responses"`
}

type Result struct {
	fx.Out
	Option health.CheckerOption `group:"health_check_options"`
}

func New(params Params) Result {
	return Result{
		Option: health.WithPeriodicCheck(
			time.Second*10,
			time.Second*1,
			NewCheck(params.DhtCrawlerActive, params.DhtServerLastResponses),
		),
	}
}
