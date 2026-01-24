package identity

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type Anon struct {
	rbac.RoleInfo
	enforcer rbac.Enforcer
}

func (a Anon) Self() Self {
	return Self{
		Permissions: slice.Map(a.Permissions, func(perm rbac.Permission) rbac.ObjectAction {
			return perm.ObjectAction()
		}),
	}
}

func (a Anon) Enforce(ctx context.Context, objectAction rbac.ObjectAction) (bool, error) {
	return a.enforcer.Enforce(
		ctx,
		rbac.SubjectRole{Role: a.Role},
		objectAction,
	)
}
