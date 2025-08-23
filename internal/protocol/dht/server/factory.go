package server

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/responder"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/socket"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// type Params struct {
// 	fx.In
// 	// Config        Config
// 	Responder     responder.Responder
// 	Logger        *zap.Logger
// 	LastResponses *atomic.Value[LastResponses]
// 	Socket        socket.Socket
// }

// type Result struct {
// 	fx.Out
// 	Server            Runner
// 	QueryDuration     prometheus.Collector `group:"prometheus_collectors"`
// 	QuerySuccessTotal prometheus.Collector `group:"prometheus_collectors"`
// 	QueryErrorTotal   prometheus.Collector `group:"prometheus_collectors"`
// 	QueryConcurrency  prometheus.Collector `group:"prometheus_collectors"`
// }

// const (
// 	namespace = "bitmagnet" // todo: Change
// )

func New(
	socket socket.Socket,
	queryTimeout QueryTimeout,
	lastResponses *atomic.Value[LastResponses],
	responder responder.Responder,
	logger *zap.Logger,
) Runner {
	// collector := newPrometheusCollector()

	// todo: Remove prometheus
	return queryLimiter{
		serverRunner: serverRunner{healthCollector{
			serverRunner: serverRunner{&server{
				socket:           socket,
				queries:          make(map[string]chan dht.RecvMsg),
				queryTimeout:     time.Duration(queryTimeout),
				responder:        responder,
				responderTimeout: time.Second * 5,
				idIssuer:         &VariantIDIssuer{},
				logger:           logger,
			}},
			lastResponses: lastResponses,
		}},
		queryLimiter: concurrency.NewKeyedLimiter(rate.Every(time.Second), 4, 1000, time.Second*20),
	}

	// return Result{
	// 	Server:            s,
	// 	QueryDuration:     collector.queryDuration,
	// 	QuerySuccessTotal: collector.querySuccessTotal,
	// 	QueryErrorTotal:   collector.queryErrorTotal,
	// 	QueryConcurrency:  collector.queryConcurrency,
	// }
}
