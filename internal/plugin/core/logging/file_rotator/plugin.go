package file_rotator

import (
	"github.com/bitmagnet-io/bitmagnet/internal/logging/filerotator"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	FileRotator *filerotator.FileRotator
}

var (
	Ref = logging.Ref.MustSub("file_rotator")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithDependencies[Config, deps](logging.Ref),
		builder.WithDefaultConfig[Config, deps](NewDefaultConfig()),
		builder.WithFxOption[Config, deps](fx.Provide(filerotator.New)),
		builder.WithWorkerRegistryOption(
			func(cfg Config, deps deps) registry.Option {
				return registry.WithWorker(
					Ref.String(),
					deps.FileRotator,
					worker.WithAutostart(),
				)
			}),
	)
)
