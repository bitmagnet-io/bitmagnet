package user

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
	"golang.org/x/crypto/bcrypt"
)

type (
	InvitationRequired bool

	EmailRequired bool

	EmailVerification bool

	PasswordMinEntropy float64

	PasswordHashingCost int

	LoginRequestsPerMinute int

	LoginRequestBurst int
)

var (
	ParamInvitationRequired = param.MustNew(
		param.Dynamic(
			param.Description[InvitationRequired]("Require invitation for registration"),
			param.Bool[InvitationRequired](),
			param.Default(InvitationRequired(true)),
		),
	)

	ParamEmailRequired = param.MustNew(
		param.Dynamic(
			param.Description[EmailRequired]("Require users to have an email address"),
			param.Bool[EmailRequired](),
			param.Default(EmailRequired(false)),
		),
	)

	ParamEmailVerification = param.MustNew(
		param.Dynamic(
			param.Description[EmailVerification]("Enable email address verification"),
			param.Bool[EmailVerification](),
			param.Default(EmailVerification(true)),
		),
	)

	ParamPasswordMinEntropy = param.MustNew(
		param.Dynamic(
			param.Description[PasswordMinEntropy]("Minimum password entropy"),
			param.Float[PasswordMinEntropy](),
			param.Min[PasswordMinEntropy](50),
			param.Default[PasswordMinEntropy](70),
		),
	)

	ParamPasswordHashingCost = param.MustNew(
		param.Dynamic(
			param.Description[PasswordHashingCost]("Cost for password hashing"),
			param.Int[PasswordHashingCost](),
			param.Min(PasswordHashingCost(bcrypt.DefaultCost)),
			param.Default(PasswordHashingCost(bcrypt.DefaultCost)),
		),
	)

	ParamLoginRequestsPerMinute = param.MustNew(
		param.Dynamic(
			param.Description[LoginRequestsPerMinute]("Login requests per minute"),
			param.Int[LoginRequestsPerMinute](),
			param.GreaterThan[LoginRequestsPerMinute](0),
			param.Default[LoginRequestsPerMinute](30),
		),
	)

	ParamLoginRequestBurst = param.MustNew(
		param.Dynamic(
			param.Description[LoginRequestBurst]("Login requests limit burst"),
			param.Int[LoginRequestBurst](),
			param.GreaterThan[LoginRequestBurst](0),
			param.Default[LoginRequestBurst](5),
		),
	)
)
