package httpserver

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	PrometheusRegistry *prometheus.Registry
}

type Result struct {
	fx.Out
	PprofOption      httpserver.Option `group:"http_server_options"`
	PrometheusOption httpserver.Option `group:"http_server_options"`
}

func New(p Params) Result {
	return Result{
		PprofOption: pprofBuilder{},
		PrometheusOption: prometheusBuilder{
			registry: p.PrometheusRegistry,
		},
	}
}
