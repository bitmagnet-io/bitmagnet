package httpserver

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Config torznab.Config
	Client lazy.Lazy[torznab.Client]
}

type Result struct {
	fx.Out
	Option httpserver.Option `group:"http_server_options"`
}

func New(p Params) Result {
	return Result{
		Option: builder{
			client: p.Client,
		},
	}
}

type builder struct {
	config torznab.Config
	client lazy.Lazy[torznab.Client]
}

func (builder) Key() string {
	return "torznab"
}

func (b builder) Apply(e *gin.Engine) error {
	client, err := b.client.Get()
	if err != nil {
		return err
	}
	h := handler{
		config: b.config,
		client: client,
	}
	e.GET("/torznab/*any", h.handleRequest)
	return nil
}
