package clientfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/client"
	"github.com/bitmagnet-io/bitmagnet/internal/config/configfx"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"client",
		configfx.NewConfigModule[client.Config]("client", client.NewDefaultConfig()),
	)
}
