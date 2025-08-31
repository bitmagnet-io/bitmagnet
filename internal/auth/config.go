package auth

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

type (
	JWTSecret   string
	JWTDuration time.Duration
)

var (
	ParamJWTSecret = param.MustNew(
		param.Description[JWTSecret]("JWT secret (if empty, a random string will be generated at runtime)"),
	)

	ParamJWTDuration = param.MustNew(
		param.Description[JWTDuration]("JWT duration"),
		param.Duration[JWTDuration](true),
		param.Default(JWTDuration(time.Hour*24)),
	)
)
