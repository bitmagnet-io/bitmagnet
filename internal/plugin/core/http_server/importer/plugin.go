package importer

import (
	"github.com/bitmagnet-io/bitmagnet/internal/importer"
	"github.com/bitmagnet-io/bitmagnet/internal/importer/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline/processor"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type deps struct {
	fx.In
	Importer importer.Importer
	Logger   *zap.Logger
}

var (
	Ref = core.Ref.MustSub("importer")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Runs an import facility at the /import endpoint"),
		builder.WithActivation[deps](plugin.ActivationEnabled),
		builder.WithDependencies[deps](
			logging.Ref,
			postgres.Ref,
			processor.Ref,
		),
		builder.WithFxOption[deps](
			fx.Provide(importer.New),
		),
		builder.WithGinOption(
			Ref,
			0,
			func(deps deps) gin.OptionFunc {
				return httpserver.New(deps.Importer, deps.Logger.Named(Ref.String()))
			},
		),
	)
)
