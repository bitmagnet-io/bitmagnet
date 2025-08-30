package httpserver

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ginContextKey string

const ginContextValue ginContextKey = "gincontext"

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginContextValue, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GinContextFromContext(ctx context.Context) (*gin.Context, bool) {
	value, ok := ctx.Value(ginContextValue).(*gin.Context)

	return value, ok
}
