package queuefx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/asynqmon"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/client"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/prometheus"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/server"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"queue",
		configfx.NewConfigModule[queue.Config]("queue", queue.NewDefaultConfig()),
		fx.Provide(
			asynqmon.New,
			client.New,
			prometheus.New,
			server.New,
		),
	)
}
