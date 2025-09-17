package rbac

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

type (
	AnonymousAccess bool

	CacheTTL time.Duration
)

var (
	ParamAnonymousAccess = param.MustNew(
		param.Dynamic(
			param.Description[AnonymousAccess]("Allow anonymous access to the application"),
			param.Bool[AnonymousAccess](),
			param.Default[AnonymousAccess](true),
		),
	)

	ParamCacheTTL = param.MustNew(
		param.Duration[CacheTTL](false),
		param.Description[CacheTTL]("Permissions cache TTL"),
		param.Default(CacheTTL(time.Minute)),
	)
)
