package httpserver

import (
	"context"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver/ginzap"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/worker"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type Params struct {
	fx.In
	Config  Config
	Options []Option `group:"http_server_options"`
	Logger  *zap.Logger
}

type Result struct {
	fx.Out
	Gin    *gin.Engine
	Server *http.Server
	Worker worker.Worker `group:"workers"`
}

func New(p Params) (r Result, err error) {
	gin.SetMode(p.Config.GinMode)
	g := gin.New()
	g.Use(ginzap.Ginzap(p.Logger.Named("gin"), time.RFC3339, true), gin.Recovery())
	for _, o := range p.Options {
		if buildErr := o.Apply(g); buildErr != nil {
			err = buildErr
			return
		}
	}
	s := &http.Server{
		Addr:    p.Config.LocalAddress,
		Handler: g.Handler(),
	}
	r.Worker = worker.NewWorker(
		"http_server",
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				ln, listenErr := net.Listen("tcp", s.Addr)
				if listenErr != nil {
					return listenErr
				}
				go (func() {
					serveErr := s.Serve(ln)
					if !errors.Is(serveErr, http.ErrServerClosed) {
						panic(serveErr)
					}
				})()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return s.Shutdown(ctx)
			},
		},
	)
	r.Gin = g
	return
}

type Option interface {
	Apply(engine *gin.Engine) error
}
