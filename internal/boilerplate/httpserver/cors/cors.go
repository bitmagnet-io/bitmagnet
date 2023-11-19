package cors

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	gincors "github.com/rs/cors/wrapper/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type Params struct {
	fx.In
	Config httpserver.Config
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Option httpserver.Option `group:"http_server_options"`
}

func New(p Params) Result {
	config := p.Config.Cors
	return Result{
		Option: corsOption{gincors.New(cors.Options{
			AllowedOrigins:       config.AllowedOrigins,
			AllowedMethods:       config.AllowedMethods,
			AllowedHeaders:       config.AllowedHeaders,
			ExposedHeaders:       config.ExposedHeaders,
			MaxAge:               config.MaxAge,
			AllowCredentials:     config.AllowCredentials,
			AllowPrivateNetwork:  config.AllowPrivateNetwork,
			OptionsPassthrough:   config.OptionsPassthrough,
			OptionsSuccessStatus: config.OptionsSuccessStatus,
			Debug:                config.Debug,
			// we don't need every request logged so apply sampling
			Logger: corsLogger{p.Logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				return zapcore.NewSamplerWithOptions(core, time.Hour, 10, 0)
			})).Named("cors")},
		})},
	}
}

type corsOption struct {
	handlerFunc gin.HandlerFunc
}

func (corsOption) Key() string {
	return "cors"
}

func (c corsOption) Apply(g *gin.Engine) error {
	g.Use(c.handlerFunc)
	return nil
}

type corsLogger struct {
	logger *zap.SugaredLogger
}

func (c corsLogger) Printf(format string, v ...interface{}) {
	c.logger.Debugf(format, v...)
}
