package identity

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type User struct {
	model.User
	role     rbac.RoleInfo
	enforcer rbac.Enforcer
}

func (u User) Self() Self {
	return Self{
		User: &u.User,
		Permissions: slice.Map(u.role.Permissions, func(perm rbac.Permission) rbac.ObjectAction {
			return perm.ObjectAction()
		}),
	}
}

func (u User) Enforce(ctx context.Context, objectAction rbac.ObjectAction) (bool, error) {
	return u.enforcer.Enforce(
		ctx,
		rbac.SubjectRole{Role: rbac.Role(u.RoleName)},
		objectAction,
	)
}
