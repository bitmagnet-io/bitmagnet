package server

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

type (
	QueryTimeout time.Duration
)

var (
	ParamQueryTimeout = param.MustNew(
		param.WithDefault(QueryTimeout(time.Second*10)),
		param.WithGreaterThan(QueryTimeout(0)),
	)
)

// type Config struct {
// 	Port          uint16
// 	QueryTimeout  time.Duration
// 	SocketAdapter string
// }

// func NewDefaultConfig() Config {
// 	return Config{
// 		Port:          3334,
// 		QueryTimeout:  time.Second * 10,
// 		SocketAdapter: "net",
// 	}
// }
