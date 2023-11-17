package server

import (
	"context"
	"fmt"
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
	Server            Server
	AppHook           fx.Hook              `group:"app_hooks"`
	QueryDuration     prometheus.Collector `group:"prometheus_collectors"`
	QuerySuccessTotal prometheus.Collector `group:"prometheus_collectors"`
	QueryErrorTotal   prometheus.Collector `group:"prometheus_collectors"`
	QueryConcurrency  prometheus.Collector `group:"prometheus_collectors"`
}

const namespace = "bitmagnet"
const subsystem = "dht_server"

func New(p Params) Result {
	s := &server{
		ready:            make(chan struct{}),
		stopped:          make(chan struct{}),
		localAddr:        netip.AddrPortFrom(netip.IPv4Unspecified(), p.Config.Port),
		socket:           NewSocket(),
		queries:          make(map[string]chan dht.RecvMsg),
		queryTimeout:     p.Config.QueryTimeout,
		responder:        p.Responder,
		responderTimeout: time.Second * 5,
		idIssuer:         &variantIdIssuer{},
		logger:           p.Logger.Named(subsystem),
	}
	collector := newPrometheusCollector(s)
	return Result{
		Server: queryLimiter{
			server:       collector,
			queryLimiter: concurrency.NewKeyedLimiter(rate.Every(time.Second), 4, 1000, time.Second*20),
		},
		AppHook: fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := s.socket.Open(s.localAddr); err != nil {
					return fmt.Errorf("could not open socket: %w", err)
				}
				go s.start()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				close(s.stopped)
				return nil
			},
		},
		QueryDuration:     collector.queryDuration,
		QuerySuccessTotal: collector.querySuccessTotal,
		QueryErrorTotal:   collector.queryErrorTotal,
		QueryConcurrency:  collector.queryConcurrency,
	}
}
