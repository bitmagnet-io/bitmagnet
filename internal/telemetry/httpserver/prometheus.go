package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewPrometheus(registry *prometheus.Registry) gin.OptionFunc {
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})

	return func(e *gin.Engine) {
		e.Any("/metrics", func(c *gin.Context) {
			h.ServeHTTP(c.Writer, c.Request)
		})
	}
}
