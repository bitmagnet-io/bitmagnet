package gqlmodel

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/user"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/auth"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type SelfQuery struct {
	RBAC rbac.Repository
	User user.Service
}

func (q *SelfQuery) Identity(ctx context.Context) (Self, error) {
	userModel, roles := auth.UserRolesFromContext(ctx)

	roleInfos, err := q.RBAC.GetRoles(ctx, roles)
	if err != nil {
		return Self{}, err
	}

	var user *User

	if userModel != nil {
		transformed := transformUser(*userModel)
		user = &transformed
	}

	return Self{
		User: user,
		Roles: slice.Map(roles, func(role rbac.Role) string {
			return string(role)
		}),
		Permissions: slice.FlatMap(roleInfos, func(roleInfo rbac.RoleInfo) []Permission {
			return slice.Map(roleInfo.Permissions, transformPermission)
		}),
	}, nil
}

func (q *SelfQuery) PasswordEntropy(password string) user.PasswordEntropyResult {
	return q.User.PasswordEntropy(password)
}

type SelfMutation struct {
	User user.Service
	RBAC rbac.Repository
}

func (m *SelfMutation) Register(ctx context.Context, input user.RegisterRequest) (RegisterResult, error) {
	usr, err := m.User.Register(ctx, input)

	if err != nil {
		return RegisterResult{}, err
	}

	return RegisterResult{
		User: transformUser(usr),
	}, nil
}

func (m *SelfMutation) Login(ctx context.Context, username, password string) (LoginResult, error) {
	success, err := m.User.Login(ctx, username, password)

	if err != nil {
		return LoginResult{}, err
	}

	roleInfos, err := m.RBAC.GetRoles(ctx, []rbac.Role{rbac.Role(success.User.RoleName)})
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		Token: success.Token,
		User:  transformUser(success.User),
		Permissions: slice.FlatMap(roleInfos, func(roleInfo rbac.RoleInfo) []Permission {
			return slice.Map(roleInfo.Permissions, transformPermission)
		}),
	}, nil
}
