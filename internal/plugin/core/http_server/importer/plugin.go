package importer

import (
	"github.com/bitmagnet-io/bitmagnet/internal/importer"
	"github.com/bitmagnet-io/bitmagnet/internal/importer/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
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
	Ref = ref.Root.MustSub("importer")

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
