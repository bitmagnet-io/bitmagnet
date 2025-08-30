package queue

import "github.com/bitmagnet-io/bitmagnet/internal/config/param"

type Autostart bool

var (
	ParamAutostart = param.MustNew(
		param.Description[Autostart]("Start the queue worker automatically"),
		param.Bool[Autostart](),
		param.Default(Autostart(true)),
	)
)
