package cache

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

type (
	TTL      time.Duration
	MaxItems int
)

var (
	ParamTTL = param.MustNew(
		param.Dynamic(
			param.Description[TTL]("TTL for the cache"),
			param.Duration[TTL](true),
			param.Default(TTL(time.Minute)),
		),
	)

	ParamMaxItems = param.MustNew(
		param.Dynamic(
			param.Description[MaxItems]("Maximum number of items to cache"),
			param.Int[MaxItems](),
			param.Default(MaxItems(1000)),
			param.GreaterThan(MaxItems(0)),
		),
	)
)
