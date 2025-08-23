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
		param.WithDynamic(
			param.WithDoc[DialTimeout]("dial timeout"),
			param.WithDefault(DialTimeout(time.Second*2)),
			param.WithGreaterThan(DialTimeout(0)),
		),
	)

	ParamRequestTimeout = param.MustNew(
		param.WithDynamic(
			param.WithDoc[RequestTimeout]("request timeout"),
			param.WithDefault(RequestTimeout(time.Second*6)),
			param.WithGreaterThan(RequestTimeout(0)),
		),
	)

	ParamMaxConcurrency = param.MustNew(
		param.WithDynamic(
			param.WithDoc[MaxConcurrency]("max concurrency"),
			param.WithDefault(MaxConcurrency(1000)),
			param.WithGreaterThan(MaxConcurrency(0)),
		),
	)
)
