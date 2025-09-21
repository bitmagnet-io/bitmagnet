package identity

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
)

type authenticatorAnon struct {
	rbac rbac.Service
}

func (a authenticatorAnon) Authenticate(ctx context.Context, token string) (Identity, bool, error) {
	roleInfo, err := a.rbac.GetRole(ctx, rbac.RoleAnon)
	if err != nil {
		return nil, true, err
	}

	return Anon{
		RoleInfo: roleInfo,
		enforcer: a.rbac,
	}, true, nil
}
