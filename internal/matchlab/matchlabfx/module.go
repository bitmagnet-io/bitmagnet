package matchlabfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/matchlab"
	"github.com/bitmagnet-io/bitmagnet/internal/matchlab/matchlabhealthcheck"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"matchlab",
		configfx.NewConfigModule[matchlab.Config]("matchlab", matchlab.NewDefaultConfig()),
		fx.Provide(matchlab.New, matchlabhealthcheck.New),
	)
}
