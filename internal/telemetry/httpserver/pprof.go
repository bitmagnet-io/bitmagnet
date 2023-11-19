package httpserver

import (
	"github.com/gin-gonic/gin"
	pyroscope "github.com/grafana/pyroscope-go/godeltaprof/http/pprof"
	"net/http/pprof"
	"runtime"
)

type pprofBuilder struct{}

func (pprofBuilder) Key() string {
	return "pprof"
}

func (b pprofBuilder) Apply(e *gin.Engine) error {
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
