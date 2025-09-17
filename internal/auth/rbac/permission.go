package rbac

import "github.com/bitmagnet-io/bitmagnet/internal/slice"

type Permission interface {
	Subject
	ObjectAction() ObjectAction
	Core() bool
	sentinel()
}

type permission struct {
	Subject
	objectAction ObjectAction
	core         bool
}

func (p permission) ObjectAction() ObjectAction {
	return p.objectAction
}

func (p permission) Core() bool {
	return p.core
}

func (permission) sentinel() {}

func NewPermission(sub Subject, objAct ObjectAction) Permission {
	return permission{
		Subject:      sub,
		objectAction: objAct,
	}
}

func CorePermissions() []Permission {
	return []Permission{
		permission{
			Subject: SubjectRole{
				Role: RoleAdmin,
			},
			objectAction: ObjectAction{
				Namespace: "**",
				Object:    "**",
				Action:    "**",
			},
		},
	}
}

type PermissionProvider func() []Permission

func PermissionProviders(providers ...PermissionProvider) PermissionProvider {
	return func() []Permission {
		return slice.FlatMap(providers, func(provider PermissionProvider) []Permission {
			return provider()
		})
	}
}
