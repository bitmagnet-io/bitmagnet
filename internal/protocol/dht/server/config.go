package server

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

type (
	QueryTimeout time.Duration
)

var ParamQueryTimeout = param.MustNew(
	param.Description[QueryTimeout]("Query timeout"),
	param.Duration[QueryTimeout](true),
	param.Default(QueryTimeout(time.Second*10)),
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
