package prometheus

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	plugin_metrics "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/telemetry/httpserver"
	"github.com/gin-gonic/gin"
	sink "github.com/hashicorp/go-metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.uber.org/fx"
)

type (
	config struct{}

	deps struct {
		fx.In
		Registry *prometheus.Registry
	}
)

var (
	Ref = plugin_metrics.Ref.MustSub("prometheus")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[config, deps](),
		builder.WithFxOption[config, deps](
			fx.Provide(
				fx.Annotate(
					func() (*prometheus.Registry, error) {
						registry := prometheus.NewRegistry()

						for _, c := range []prometheus.Collector{
							collectors.NewGoCollector(),
							collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
								Namespace: "bitmagnet",
							}),
						} {
							if err := registry.Register(c); err != nil {
								return nil, err
							}
						}

						return registry, nil
					},
				),
				func(registry *prometheus.Registry) (*sink.PrometheusSink, error) {
					return sink.NewPrometheusSinkFrom(
						sink.PrometheusOpts{
							Expiration: 60 * time.Minute,
							Registerer: registry,
						},
					)
				},
				fx.Annotate(
					func(sink *sink.PrometheusSink) metrics.Option {
						return metrics.WithSink(sink)
					},
					fx.ResultTags(`group:"metrics_options"`),
				),
			),
		),
		builder.WithGinOption(
			Ref.MustSub("pprof"),
			func(config, deps) gin.OptionFunc {
				return httpserver.NewPProf()
			},
		),
		// todo: Move this
		builder.WithGinOption(
			Ref,
			func(_ config, deps deps) gin.OptionFunc {
				return httpserver.NewPrometheus(deps.Registry)
			},
		),
	)
)
