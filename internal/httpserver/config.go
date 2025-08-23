package httpserver

import "github.com/bitmagnet-io/bitmagnet/internal/config/param"

type LocalAddress string

var ParamLocalAddress = param.MustNew(
	param.WithDefault(LocalAddress(":3333")),
)
