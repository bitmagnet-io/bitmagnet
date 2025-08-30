package httpserver

import "github.com/bitmagnet-io/bitmagnet/internal/config/param"

type LocalAddress string

var ParamLocalAddress = param.MustNew(
	param.Description[LocalAddress]("Local address on which to listen"),
	param.Default(LocalAddress(":3333")),
)
