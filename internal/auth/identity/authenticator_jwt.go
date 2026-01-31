package identity

import (
	"context"
	"errors"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/jwt"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/user"
)

type authenticatorJWT struct {
	jwtService  jwt.Service
	userService user.Service
	rbac        rbac.Service
}

func (a authenticatorJWT) Authenticate(ctx context.Context, token string) (Identity, bool, error) {
	claims, err := a.jwtService.Parse(token)
	if err != nil {
		return nil, errors.Is(err, jwt.ErrTokenInvalidClaims), err
	}

	user, err := a.userService.Get(ctx, claims.UserID)
	if err != nil {
		return nil, true, err
	}

	role, err := a.rbac.GetRole(ctx, rbac.Role(user.RoleName))
	if err != nil {
		return nil, true, err
	}

	return User{
		User:     user,
		role:     role,
		enforcer: a.rbac,
	}, true, nil
}
