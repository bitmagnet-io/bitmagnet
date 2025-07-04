package ginzap

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Option(logger *zap.Logger) httpserver.Option {
	return httpserver.NewOption("logger", func(engine *gin.Engine) {
		engine.Use(Ginzap(logger, time.RFC3339, true))
	})
}
