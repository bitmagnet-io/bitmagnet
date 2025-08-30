package file_rotator

import (
	"github.com/bitmagnet-io/bitmagnet/internal/fs"
	"github.com/bitmagnet-io/bitmagnet/internal/logging/file_rotator"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	FileRotator *file_rotator.FileRotator
}

var (
	Ref = logging.Ref.MustSub("file_rotator")
	// todo: Where's the zap core?
	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Outputs JSON logs to rotated log files"),
		builder.WithActivation[deps](plugin.ActivationDisabled),
		builder.WithDependencies[deps](logging.Ref),
		builder.WithConfig[deps](Ref.MustSub("level"), file_rotator.ParamLevel),
		builder.WithConfig[deps](Ref.MustSub("sub_path"), file_rotator.ParamSubPath),
		builder.WithConfig[deps](Ref.MustSub("base_name"), file_rotator.ParamBaseName),
		builder.WithConfig[deps](Ref.MustSub("max_age"), file_rotator.ParamMaxAge),
		builder.WithConfig[deps](Ref.MustSub("max_size"), file_rotator.ParamMaxSize),
		builder.WithConfig[deps](Ref.MustSub("max_backups"), file_rotator.ParamMaxBackups),
		builder.WithConfig[deps](Ref.MustSub("buffer_size"), file_rotator.ParamBufferSize),
		builder.WithFxOption[deps](
			fx.Provide(
				func(fsProvider fs.FSProvider) file_rotator.FS {
					return fsProvider.FSData()
				},
				file_rotator.New,
			),
		),
		builder.WithWorker(
			func(deps deps) (runner.Provider, worker.Option) {
				return deps.FileRotator,
					worker.WithAutostart(true)
			},
		),
	)
)
