package info_hash_blocker

import (
	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/migrator"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	Blocker blocker.Blocker
}

var (
	Ref = database.Ref.MustSub("info_hash_blocker")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps](
			"Maintains a stable bloom filter of info hashes that have been deleted and blocked, preventing them from being crawled again",
		),
		builder.WithDependencies[deps](
			migrator.Ref,
			postgres.Ref,
		),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithFxOption[deps](
			fx.Provide(blocker.New),
		),
		builder.WithWorker(
			func(deps deps) (runner.Provider, worker.Option) {
				return deps.Blocker, worker.WithDependencies(
					migrator.Ref,
					postgres.Ref,
				)
			},
		),
	)
)
