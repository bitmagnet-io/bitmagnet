package httpserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver/ginzap"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/worker"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
	"sort"
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
	options, optionsErr := resolveOptions(p.Config.Options, p.Options)
	if optionsErr != nil {
		err = optionsErr
		return
	}
	for _, o := range options {
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

func resolveOptions(param []string, options []Option) ([]Option, error) {
	paramMap := make(map[string]struct{})
	for _, p := range param {
		paramMap[p] = struct{}{}
	}
	enabledOptions := make([]Option, 0, len(options))
	foundMap := make(map[string]struct{}, len(options))
	for _, o := range options {
		if _, ok := foundMap[o.Key()]; ok {
			return nil, fmt.Errorf("duplicate http server option: '%s'", o.Key())
		}
		foundMap[o.Key()] = struct{}{}
		enabled := false
		if _, ok := paramMap["*"]; ok {
			enabled = true
		} else if _, ok := paramMap[o.Key()]; ok {
			enabled = true
		}
		if enabled {
			enabledOptions = append(enabledOptions, o)
		}
	}
	for p := range paramMap {
		if _, ok := foundMap[p]; !ok && p != "*" {
			return nil, fmt.Errorf("unknown http server option: '%s'", p)
		}
	}
	sort.Slice(enabledOptions, func(i, j int) bool {
		return enabledOptions[i].Key() < enabledOptions[j].Key()
	})
	return enabledOptions, nil
}

type Option interface {
	Key() string
	Apply(engine *gin.Engine) error
}
