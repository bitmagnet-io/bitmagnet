package persister

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

type (
	MaxSize int
	MaxWait int64
)

var (
	ParamMaxSize = param.MustNew(
		param.Dynamic(
			param.Description[MaxSize]("Maximum buffer size"),
			param.Int[MaxSize](),
			param.Default(MaxSize(1000)),
			param.GreaterThan(MaxSize(0)),
		),
	)

	ParamMaxWait = param.MustNew(
		param.Dynamic(
			param.Description[MaxWait]("Maximum time to buffer items before flushing"),
			param.Duration[MaxWait](true),
			param.Default(MaxWait(time.Second*10)),
		),
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
