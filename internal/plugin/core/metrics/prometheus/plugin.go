package prometheus

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/telemetry/httpserver"
	"github.com/gin-gonic/gin"
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
	Ref = metrics.Ref.MustSub("prometheus")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[config, deps](),
		builder.WithFxOption[config, deps](
			fx.Provide(
				fx.Annotate(
					func(cs []prometheus.Collector) (*prometheus.Registry, error) {
						cs = append(
							[]prometheus.Collector{
								collectors.NewGoCollector(),
								collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
									Namespace: "bitmagnet",
								}),
							}, cs...)

						registry := prometheus.NewRegistry()

						for _, c := range cs {
							if err := registry.Register(c); err != nil {
								return nil, err
							}
						}

						return registry, nil
					},
					fx.ParamTags(`group:"prometheus_collectors"`),
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
