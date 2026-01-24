package socket

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

type (
	Port        uint16
	AdapterName string
)

var ParamPort = param.MustNew(
	param.Description[Port]("Socket port number"),
	param.PortNumber[Port](),
	param.Default(Port(3334)),
)

func ParamAdapter() param.Param[AdapterName] {
	return param.MustNew(
		param.Description[AdapterName]("Socket adapter name"),
		param.Default(AdapterName("net")),
		param.EnumValues(AdapterNames()...),
	)
}

// type Config struct {
// 	Adapter   string
// 	LocalAddr netip.AddrPort
// }
