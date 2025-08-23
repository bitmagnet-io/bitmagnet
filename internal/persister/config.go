package persister

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

type (
	MaxSize int
	MaxWait time.Duration
)

var (
	ParamMaxSize = param.MustNew(
		param.WithDefault(MaxSize(1000)),
		param.WithGreaterThan(MaxSize(0)),
	)

	ParamMaxWait = param.MustNew(
		param.WithDefault(MaxWait(time.Second*10)),
		param.WithGreaterThan(MaxWait(0)),
	)
)

// type Config struct {
// 	MaxSize int
// 	MaxWait time.Duration
// }

// func NewDefaultConfig() Config {
// 	return Config{
// 		MaxSize: 1000,
// 		MaxWait: 10 * time.Second,
// 	}
// }
