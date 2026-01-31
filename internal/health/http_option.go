package health

import (
	"github.com/gin-gonic/gin"
)

func NewHTTPOption(checker Checker) gin.OptionFunc {
	return func(e *gin.Engine) {
		handler := NewHandler(checker)

		e.GET("/status", func(c *gin.Context) {
			handler(c.Writer, c.Request)
		})
	}
}
