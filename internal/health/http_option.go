package health

import (
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/gin-gonic/gin"
)

func NewHTTPOption(checker Checker) httpserver.Option {
	return handlerBuilder{checker}
}

type handlerBuilder struct {
	checker Checker
}

func (handlerBuilder) Key() string {
	return "health"
}

func (b handlerBuilder) Apply(e *gin.Engine) {
	handler := NewHandler(b.checker)

	e.GET("/status", func(c *gin.Context) {
		handler(c.Writer, c.Request)
	})
}
