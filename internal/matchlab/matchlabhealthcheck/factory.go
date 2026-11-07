package matchlabhealthcheck

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/matchlab"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Config matchlab.Config
	Client lazy.Lazy[matchlab.Client]
}
type Result struct {
	fx.Out
	Option health.CheckerOption `group:"health_check_options"`
}

func New(p Params) Result {
	return Result{Option: health.WithPeriodicCheck(
		5*time.Minute, 0, NewCheck(p.Config, p.Client),
	)}
}
