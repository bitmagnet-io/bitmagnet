package server

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/responder"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/socket"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type Params struct {
	fx.In
	Config        Config
	Responder     responder.Responder
	Logger        *zap.Logger
	LastResponses *concurrency.AtomicValue[LastResponses]
	Socket        socket.Socket
}

type Result struct {
	fx.Out
	Server            Runner
	QueryDuration     prometheus.Collector `group:"prometheus_collectors"`
	QuerySuccessTotal prometheus.Collector `group:"prometheus_collectors"`
	QueryErrorTotal   prometheus.Collector `group:"prometheus_collectors"`
	QueryConcurrency  prometheus.Collector `group:"prometheus_collectors"`
}

const (
	namespace = "bitmagnet" // todo: Change
)

func New(p Params) Result {
	collector := newPrometheusCollector()

	s := queryLimiter{
		serverRunner: serverRunner{prometheusServerWrapper{
			prometheusCollector: collector,
			serverRunner: serverRunner{healthCollector{
				serverRunner: serverRunner{&server{
					socket:           p.Socket,
					queries:          make(map[string]chan dht.RecvMsg),
					queryTimeout:     p.Config.QueryTimeout,
					responder:        p.Responder,
					responderTimeout: time.Second * 5,
					idIssuer:         &VariantIDIssuer{},
					logger:           p.Logger.Named(Namespace),
				}},
				lastResponses: p.LastResponses,
			},
			}}},
		queryLimiter: concurrency.NewKeyedLimiter(rate.Every(time.Second), 4, 1000, time.Second*20),
	}

	return Result{
		Server:            s,
		QueryDuration:     collector.queryDuration,
		QuerySuccessTotal: collector.querySuccessTotal,
		QueryErrorTotal:   collector.queryErrorTotal,
		QueryConcurrency:  collector.queryConcurrency,
	}
}
