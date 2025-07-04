package httpserver

import (
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"github.com/gin-gonic/gin"
)

func New(client torznab.Client, config torznab.Config) httpserver.Option {
	return builder{
		client: client,
		config: config,
	}
}

type builder struct {
	client torznab.Client
	config torznab.Config
}

func (builder) Key() string {
	return "torznab"
}

func (b builder) Apply(e *gin.Engine) {
	h := handler{
		config: b.config,
		client: b.client,
	}
	e.GET("/torznab/*any", h.handleRequest)
}
