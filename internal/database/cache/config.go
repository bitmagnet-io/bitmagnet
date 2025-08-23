package cache

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

type (
	TTL     time.Duration
	MaxKeys int
)

var (
	ParamTTL = param.MustNew(
		param.WithDynamic(
			param.WithDefault(TTL(time.Minute)),
		),
	)

	ParamMaxKeys = param.MustNew(
		param.WithDynamic(
			param.WithDefault(MaxKeys(1000)),
			param.WithGreaterThan(MaxKeys(0)),
		),
	)
)
