package httpserver

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver"
	"github.com/gin-gonic/gin"
	"github.com/hellofresh/health-go/v5"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Health *health.Health
}

type Result struct {
	fx.Out
	Option httpserver.Option `group:"http_server_options"`
}

func New(p Params) (r Result, err error) {
	handler := p.Health.Handler()
	r.Option = &builder{
		handler: func(c *gin.Context) {
			handler.ServeHTTP(c.Writer, c.Request)
		},
	}
	return
}

type builder struct {
	handler gin.HandlerFunc
}

func (builder) Key() string {
	return "status"
}

func (b *builder) Apply(e *gin.Engine) error {
	e.GET("/status", b.handler)
	return nil
}
