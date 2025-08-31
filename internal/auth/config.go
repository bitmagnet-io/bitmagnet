package auth

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

type (
	JWTSecret            string
	JWTDuration          time.Duration
	PasswordMinLength    int
	PasswordMaxLength    int
	PasswordMinUpper     int
	PasswordMinLower     int
	PasswordMinDigit     int
	PasswordMinSpecial   int
	PasswordSpecialChars string
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

	ParamPasswordMinLength = param.MustNew(
		param.Description[PasswordMinLength]("Minimum password length"),
		param.Int[PasswordMinLength](),
		param.Default(PasswordMinLength(8)),
		param.Min(PasswordMinLength(4)),
		param.Max(PasswordMinLength(72)),
	)

	ParamPasswordMaxLength = param.MustNew(
		param.Description[PasswordMaxLength]("Maximum password length"),
		param.Int[PasswordMaxLength](),
		param.Default(PasswordMaxLength(72)), // maximum length for bcrypt hashing
		param.Min(PasswordMaxLength(4)),
		param.Max(PasswordMaxLength(72)),
	)

	ParamPasswordMinUpper = param.MustNew(
		param.Description[PasswordMinUpper]("Minimum number of uppercase letters required"),
		param.Int[PasswordMinUpper](),
		param.Default(PasswordMinUpper(1)),
		param.Min(PasswordMinUpper(0)),
	)

	ParamPasswordMinLower = param.MustNew(
		param.Description[PasswordMinLower]("Minimum number of lowercase letters required"),
		param.Int[PasswordMinLower](),
		param.Default(PasswordMinLower(1)),
		param.Min(PasswordMinLower(0)),
	)

	ParamPasswordMinDigit = param.MustNew(
		param.Description[PasswordMinDigit]("Minimum number of digits required"),
		param.Int[PasswordMinDigit](),
		param.Default(PasswordMinDigit(1)),
		param.Min(PasswordMinDigit(0)),
	)

	ParamPasswordMinSpecial = param.MustNew(
		param.Description[PasswordMinSpecial]("Minimum number of special characters required"),
		param.Int[PasswordMinSpecial](),
		param.Default(PasswordMinSpecial(0)),
		param.Min(PasswordMinSpecial(0)),
	)

	ParamPasswordSpecialChars = param.MustNew(
		param.Description[PasswordSpecialChars]("Special characters allowed in passwords"),
		param.Default(PasswordSpecialChars("!@#$%^&*()_+-=[]{}|;:,.<>?")),
	)
)

type PasswordPolicyConfig struct {
	MinLength    PasswordMinLength
	MaxLength    PasswordMaxLength
	MinUpper     PasswordMinUpper
	MinLower     PasswordMinLower
	MinDigit     PasswordMinDigit
	MinSpecial   PasswordMinSpecial
	SpecialChars PasswordSpecialChars
}
