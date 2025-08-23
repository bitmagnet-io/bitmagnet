package socket

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

type (
	Port        int
	AdapterName string
)

var (
	ParamPort = param.MustNew(
		param.WithDefault(Port(3334)),
		param.WithMin(Port(0)),
		param.WithMax(Port(65535)),
	)

	ParamAdapter = param.MustNew(
		param.WithDefault(AdapterName("net")),
		param.WithEnumValues(AdapterNames()...),
	)
)

// type Config struct {
// 	Adapter   string
// 	LocalAddr netip.AddrPort
// }
