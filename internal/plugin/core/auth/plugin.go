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
	Logger      *zap.Logger
	JWTSecret   auth.JWTSecret
	JWTDuration auth.JWTDuration
	AuthService auth.AuthService
	JWTService  auth.JWTService
}

var (
	Ref = core.Ref.MustSub("auth")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides user authentication services"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithConfig[deps](Ref.MustSub("jwt_secret"), auth.ParamJWTSecret),
		builder.WithConfig[deps](Ref.MustSub("jwt_duration"), auth.ParamJWTDuration),
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
