package asynqmon

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/redis"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynqmon"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Redis *redis.Client
}

type Result struct {
	fx.Out
	HttpOption httpserver.Option `group:"http_server_options"`
}

const rootPath = "/asynqmon"

func New(p Params) Result {
	return Result{
		HttpOption: asynqmonBuilder{
			options: asynqmon.Options{
				RootPath:     rootPath, // RootPath specifies the root for asynqmon app
				RedisConnOpt: redis.Wrapper{Redis: p.Redis},
			},
		},
	}
}

type asynqmonBuilder struct {
	options asynqmon.Options
}

func (b asynqmonBuilder) Apply(e *gin.Engine) error {
	handler := asynqmon.New(b.options)
	e.Any(rootPath+"/*path", func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	})
	return nil
}
