package adultfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/adult/tpdb"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"adult",
		configfx.NewConfigModule[tpdb.Config]("tpdb", tpdb.NewDefaultConfig()),
		fx.Provide(
			tpdb.New,
		),
	)
}
