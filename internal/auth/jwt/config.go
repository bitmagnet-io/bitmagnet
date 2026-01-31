package jwt

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

type (
	Secret   string
	Duration time.Duration
)

var (
	ParamSecret = param.MustNew(
		param.Description[Secret]("JWT secret (if empty, a random string will be generated at runtime)"),
	)

	ParamDuration = param.MustNew(
		param.Description[Duration]("JWT validity duration"),
		param.Duration[Duration](true),
		param.Default(Duration(time.Hour*24)),
	)
)
