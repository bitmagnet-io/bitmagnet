package metainforequester

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

type (
	DialTimeout    time.Duration
	RequestTimeout time.Duration
	MaxConcurrency int
)

var (
	ParamDialTimeout = param.MustNew(
		param.Dynamic(
			param.Description[DialTimeout]("Dial timeout"),
			param.Duration[DialTimeout](true),
			param.Default(DialTimeout(time.Second*2)),
		),
	)

	ParamRequestTimeout = param.MustNew(
		param.Dynamic(
			param.Description[RequestTimeout]("Request timeout"),
			param.Duration[RequestTimeout](true),
			param.Default(RequestTimeout(time.Second*6)),
		),
	)

	ParamMaxConcurrency = param.MustNew(
		param.Dynamic(
			param.Description[MaxConcurrency]("Maximum request concurrency"),
			param.Int[MaxConcurrency](),
			param.Default(MaxConcurrency(100)),
			param.GreaterThan(MaxConcurrency(0)),
		),
	)
)
