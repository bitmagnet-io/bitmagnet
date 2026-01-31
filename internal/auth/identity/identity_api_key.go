package identity

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type APIKey struct {
	model.APIKey
	anon     rbac.RoleInfo
	enforcer rbac.Enforcer
}

func (a APIKey) Self() Self {
	return Self{
		User:   &a.User,
		APIKey: &a.APIKey,
		Permissions: append(slice.Map(a.Permissions, func(perm model.APIKeyPermission) rbac.ObjectAction {
			return rbac.ObjectAction{
				Namespace: perm.Namespace,
				Object:    perm.Object,
				Action:    perm.Action,
			}
		}), slice.Map(a.anon.Permissions, func(perm rbac.Permission) rbac.ObjectAction {
			return perm.ObjectAction()
		})...),
	}
}

func (a APIKey) Enforce(ctx context.Context, objectAction rbac.ObjectAction) (bool, error) {
	if allow, err := a.enforcer.Enforce(
		ctx,
		rbac.SubjectRole{Role: rbac.Role(a.User.RoleName)},
		objectAction,
	); err != nil || !allow {
		return false, err
	}

	return a.enforcer.EnforceAny(
		ctx,
		append(slice.Map(a.Permissions, func(perm model.APIKeyPermission) rbac.Subject {
			return rbac.SubjectPermission{
				ObjectAction: rbac.ObjectAction{
					Namespace: perm.Namespace,
					Object:    perm.Object,
					Action:    perm.Action,
				},
			}
		}), rbac.SubjectRole{Role: a.anon.Role}),
		objectAction,
	)
}
