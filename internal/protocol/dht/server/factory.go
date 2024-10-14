package server

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/responder"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"net/netip"
	"time"
)

type Params struct {
	fx.In
	Config    Config
	Responder responder.Responder
	Logger    *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Server            lazy.Lazy[Server]
	LastResponses     *concurrency.AtomicValue[LastResponses] `name:"dht_server_last_responses"`
	AppHook           fx.Hook                                 `group:"app_hooks"`
	QueryDuration     prometheus.Collector                    `group:"prometheus_collectors"`
	QuerySuccessTotal prometheus.Collector                    `group:"prometheus_collectors"`
	QueryErrorTotal   prometheus.Collector                    `group:"prometheus_collectors"`
	QueryConcurrency  prometheus.Collector                    `group:"prometheus_collectors"`
}

const namespace = "bitmagnet"
const subsystem = "dht_server"

func New(p Params) Result {
	lastResponses := &concurrency.AtomicValue[LastResponses]{}
	collector := newPrometheusCollector()
	ls := lazy.New(func() (Server, error) {
		s := queryLimiter{
			server: prometheusServerWrapper{
				prometheusCollector: collector,
				server: healthCollector{
					baseServer: &server{
						stopped:          make(chan struct{}),
						localAddr:        netip.AddrPortFrom(netip.IPv4Unspecified(), p.Config.Port),
						socket:           NewSocket(),
						queries:          make(map[string]chan dht.RecvMsg),
						queryTimeout:     p.Config.QueryTimeout,
						responder:        p.Responder,
						responderTimeout: time.Second * 5,
						idIssuer:         &variantIdIssuer{},
						logger:           p.Logger.Named(subsystem),
					},
					lastResponses: lastResponses,
				},
			},
			queryLimiter: concurrency.NewKeyedLimiter(rate.Every(time.Second), 4, 1000, time.Second*20),
		}
		if err := s.start(); err != nil {
			return nil, fmt.Errorf("could not start server: %w", err)
		}
		return s, nil
	})
	return Result{
		Server: ls,
		AppHook: fx.Hook{
			OnStop: func(context.Context) error {
				return ls.IfInitialized(func(s Server) error {
					s.stop()
					return nil
				})
			},
		},
		LastResponses:     lastResponses,
		QueryDuration:     collector.queryDuration,
		QuerySuccessTotal: collector.querySuccessTotal,
		QueryErrorTotal:   collector.queryErrorTotal,
		QueryConcurrency:  collector.queryConcurrency,
	}
}
