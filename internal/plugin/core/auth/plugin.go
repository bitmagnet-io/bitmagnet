package auth

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/bitmagnet-io/bitmagnet/internal/auth"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type deps struct {
	fx.In
	Logger                *zap.Logger
	JWTSecret             auth.JWTSecret
	JWTDuration           auth.JWTDuration
	PasswordPolicyConfig  auth.PasswordPolicyConfig
	AuthService           auth.AuthService
	JWTService            auth.JWTService
	PasswordPolicyService auth.PasswordPolicyService
}

var (
	Ref = core.Ref.MustSub("auth")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides user authentication services"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithConfig[deps](Ref.MustSub("jwt_secret"), auth.ParamJWTSecret),
		builder.WithConfig[deps](Ref.MustSub("jwt_duration"), auth.ParamJWTDuration),
		builder.WithConfig[deps](Ref.MustSub("password_min_length"), auth.ParamPasswordMinLength),
		builder.WithConfig[deps](Ref.MustSub("password_max_length"), auth.ParamPasswordMaxLength),
		builder.WithConfig[deps](Ref.MustSub("password_min_upper"), auth.ParamPasswordMinUpper),
		builder.WithConfig[deps](Ref.MustSub("password_min_lower"), auth.ParamPasswordMinLower),
		builder.WithConfig[deps](Ref.MustSub("password_min_digit"), auth.ParamPasswordMinDigit),
		builder.WithConfig[deps](Ref.MustSub("password_min_special"), auth.ParamPasswordMinSpecial),
		builder.WithConfig[deps](Ref.MustSub("password_special_chars"), auth.ParamPasswordSpecialChars),
		builder.WithFxOption[deps](
			fx.Decorate(
				func(secret auth.JWTSecret) auth.JWTSecret {
					if secret == "" {
						secret = auth.JWTSecret(generateRandomString(32))
					}

					return secret
				},
			),
			fx.Provide(
				func(
					minLength auth.PasswordMinLength,
					maxLength auth.PasswordMaxLength,
					minUpper auth.PasswordMinUpper,
					minLower auth.PasswordMinLower,
					minDigit auth.PasswordMinDigit,
					minSpecial auth.PasswordMinSpecial,
					specialChars auth.PasswordSpecialChars,
				) auth.PasswordPolicyConfig {
					return auth.PasswordPolicyConfig{
						MinLength:    minLength,
						MaxLength:    maxLength,
						MinUpper:     minUpper,
						MinLower:     minLower,
						MinDigit:     minDigit,
						MinSpecial:   minSpecial,
						SpecialChars: specialChars,
					}
				},
				auth.NewPasswordPolicyService,
				auth.NewAuthService,
				auth.NewJWTService,
				auth.NewAuthMiddleware,
			),
		),
	)
)

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)

	return hex.EncodeToString(bytes)
}
