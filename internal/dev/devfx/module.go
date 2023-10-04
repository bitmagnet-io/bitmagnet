package devfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/boilerplatefx"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/migrations"
	"github.com/bitmagnet-io/bitmagnet/internal/database/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/dev/app/cmd/gormcmd"
	"github.com/bitmagnet-io/bitmagnet/internal/dev/app/cmd/migratecmd"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"dev",
		boilerplatefx.New(),
		configfx.NewConfigModule[postgres.Config]("postgres", postgres.NewDefaultConfig()),
		fx.Provide(database.New),
		fx.Provide(migrations.New),
		fx.Provide(postgres.New),
		fx.Provide(gormcmd.New),
		fx.Provide(migratecmd.New),
	)
}
