package matchlab

import (
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Config Config
}
type Result struct {
	fx.Out
	Client lazy.Lazy[Client]
}

func New(p Params) Result {
	return Result{Client: lazy.New(func() (Client, error) {
		return NewClient(p.Config)
	})}
}
