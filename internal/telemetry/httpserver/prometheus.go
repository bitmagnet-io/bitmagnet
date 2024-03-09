package httpserver

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusBuilder struct {
	registry lazy.Lazy[*prometheus.Registry]
}

func (prometheusBuilder) Key() string {
	return "prometheus"
}

func (b prometheusBuilder) Apply(e *gin.Engine) error {
	r, err := b.registry.Get()
	if err != nil {
		return err
	}
	h := promhttp.HandlerFor(r, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})
	e.Any("/metrics", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})
	return nil
}
