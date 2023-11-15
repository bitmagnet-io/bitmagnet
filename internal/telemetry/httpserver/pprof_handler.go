package httpserver

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver"
	"github.com/gin-gonic/gin"
	pyroscope "github.com/grafana/pyroscope-go/godeltaprof/http/pprof"
	"go.uber.org/fx"
	"net/http/pprof"
	"runtime"
)

type Result struct {
	fx.Out
	Option httpserver.Option `group:"http_server_options"`
}

func New() Result {
	return Result{
		Option: &builder{},
	}
}

type builder struct{}

func (b *builder) Priority() int {
	return 100
}

func (b *builder) Apply(e *gin.Engine) error {
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)
	e.Any("/debug/pprof/", func(c *gin.Context) {
		pprof.Index(c.Writer, c.Request)
	})
	e.Any("/debug/pprof/cmdline", func(c *gin.Context) {
		pprof.Cmdline(c.Writer, c.Request)
	})
	e.Any("/debug/pprof/profile", func(c *gin.Context) {
		pprof.Profile(c.Writer, c.Request)
	})
	e.Any("/debug/pprof/symbol", func(c *gin.Context) {
		pprof.Symbol(c.Writer, c.Request)
	})
	e.Any("/debug/pprof/trace", func(c *gin.Context) {
		pprof.Trace(c.Writer, c.Request)
	})
	e.Any("/debug/pprof/delta_heap", func(c *gin.Context) {
		pyroscope.Heap(c.Writer, c.Request)
	})
	e.Any("/debug/pprof/delta_block", func(c *gin.Context) {
		pyroscope.Block(c.Writer, c.Request)
	})
	e.Any("/debug/pprof/delta_mutex", func(c *gin.Context) {
		pyroscope.Mutex(c.Writer, c.Request)
	})
	e.Any("/debug/pprof/:profile", func(c *gin.Context) {
		pprof.Handler(c.Param("profile")).ServeHTTP(c.Writer, c.Request)
	})
	return nil
}
