package queuefx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/client"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/server"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"subscriber",
		configfx.NewConfigModule[queue.Config]("queue", queue.NewDefaultConfig()),
		fx.Provide(
			client.New,
			server.New,
		),
	)
}
