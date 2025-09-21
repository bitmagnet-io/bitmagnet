package identity

import (
	"github.com/bitmagnet-io/bitmagnet/internal/auth/api_key"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/jwt"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/user"
)

func NewAuthenticator(
	jwtService jwt.Service,
	userService user.Service,
	apiKeyService api_key.Service,
	rbacService rbac.Service,
) Authenticator {
	return authenticatorChain{
		authenticatorJWT{
			jwtService:  jwtService,
			userService: userService,
			rbac:        rbacService,
		},
		authenticatorAPIKey{
			apiKeyService: apiKeyService,
			rbac:          rbacService,
		},
		authenticatorAnon{
			rbac: rbacService,
		},
	}
}
