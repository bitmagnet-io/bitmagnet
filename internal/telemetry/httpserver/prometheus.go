package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusBuilder struct {
	registry *prometheus.Registry
}

func (b prometheusBuilder) Priority() int {
	return 0
}

func (b prometheusBuilder) Apply(e *gin.Engine) error {
	h := promhttp.HandlerFor(b.registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})
	e.Any("/metrics", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})
	return nil
}
