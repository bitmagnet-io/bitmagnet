package responder

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
)

type Params struct {
	fx.In
	KTableFactory   func() ktable.Table
	DiscoveredNodes concurrency.BatchingChannel[ktable.Node] `name:"dht_discovered_nodes"`
	Logger          *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Responder         Responder
	QueryDuration     prometheus.Collector `group:"prometheus_collectors"`
	QuerySuccessTotal prometheus.Collector `group:"prometheus_collectors"`
	QueryErrorTotal   prometheus.Collector `group:"prometheus_collectors"`
	QueryConcurrency  prometheus.Collector `group:"prometheus_collectors"`
}

const (
	namespace = "bitmagnet"
	subsystem = "dht_responder"
)

func New(p Params) Result {
	v4table := p.KTableFactory()
	v6table := p.KTableFactory()
	if v4table.Origin() != v6table.Origin() {
		panic("v4 and v6 tables have different origins")
	}
	collector := newPrometheusCollector(responderLimiter{
		responder: responder{
			nodeID:                   v4table.Origin(),
			kTable:                   v4table,
			kTable6:                  v6table,
			tokenSecret:              protocol.RandomNodeID().Bytes(),
			sampleInfoHashesInterval: 10,
		},
		limiter: NewLimiter(rate.Every(time.Second/50), 20, rate.Every(time.Second), 10, 1000, time.Second*20),
	})

	return Result{
		Responder: responderNodeDiscovery{
			responder: responderLogger{
				responder: collector,
				logger: p.Logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
					return zapcore.NewSamplerWithOptions(core, time.Minute, 10, 0)
				})).Named(subsystem),
			},
			discoveredNodes: p.DiscoveredNodes.In(),
		},
		QueryDuration:     collector.queryDuration,
		QuerySuccessTotal: collector.querySuccessTotal,
		QueryErrorTotal:   collector.queryErrorTotal,
		QueryConcurrency:  collector.queryConcurrency,
	}
}
