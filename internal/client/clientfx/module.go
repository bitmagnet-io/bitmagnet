package clientfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/client"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"client",
		configfx.NewConfigModule[client.Config]("client", client.NewDefaultConfig()),
	)
}
