package gqlmodel

import (
	"context"
	"errors"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/user"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/auth"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type AuthQuery struct {
	RBAC rbac.Service
	User user.Service
}

type User struct {
	ID          int
	Username    string
	Role        string
	Email       *string
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func transformUser(user model.User) User {
	var email *string
	if user.Email.Valid {
		email = &user.Email.String
	}

	var lastLoginAt *time.Time
	if user.LastLoginAt.Valid {
		lastLoginAt = &user.LastLoginAt.Time
	}

	return User{
		ID:          user.ID,
		Username:    user.Username,
		Role:        user.RoleName,
		Email:       email,
		LastLoginAt: lastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

type ListUsersResult struct {
	Users      []User
	TotalCount int
}

type Role struct {
	Name        string
	Core        bool
	Permissions []Permission
}

func transformRole(role rbac.RoleInfo) Role {
	return Role{
		Name: string(role.Role),
		Core: role.Core,
		Permissions: slice.Map(role.Permissions, func(perm rbac.Permission) Permission {
			return Permission{
				Subject: AuthSubject{
					Type: gen.AuthSubjectTypeRole,
					Name: string(role.Role),
				},
				ObjectAction: perm.ObjectAction(),
				Core:         perm.Core(),
			}
		}),
	}
}

func transformPermission(perm rbac.Permission) Permission {
	var subjectType gen.AuthSubjectType

	switch perm.SubjectType() {
	case rbac.SubjectTypeRole:
		subjectType = gen.AuthSubjectTypeRole
	}

	return Permission{
		Subject: AuthSubject{
			Type: subjectType,
			Name: perm.SubjectName(),
		},
		ObjectAction: perm.ObjectAction(),
		Core:         perm.Core(),
	}
}

type LoginResult struct {
	Token       string
	User        User
	Permissions []Permission
}

type RegisterResult struct {
	User User
}

type AuthSubject struct {
	Type gen.AuthSubjectType
	Name string
}

type AuthObjectAction = rbac.ObjectAction

type Permission struct {
	Subject      AuthSubject
	ObjectAction AuthObjectAction
	Core         bool
}

type Self struct {
	User        *User
	Roles       []string
	Permissions []Permission
}

func (q *AuthQuery) ListUsers(ctx context.Context, input gen.ListUsersInput) (ListUsersResult, error) {
	var params user.ListUsersParams

	if pagination, ok := input.Pagination.ValueOK(); ok {
		if page, ok := pagination.Page.ValueOK(); ok {
			params.Page = *page
		}

		if offset, ok := pagination.Offset.ValueOK(); ok {
			params.Offset = *offset
		}

		if limit, ok := pagination.Limit.ValueOK(); ok {
			params.Limit = *limit
		}
	}

	if usernameLike, ok := input.UsernameLike.ValueOK(); ok {
		params.UsernameLike = *usernameLike
	}

	result, err := q.User.List(ctx, params)
	if err != nil {
		return ListUsersResult{}, err
	}

	return ListUsersResult{
		Users:      slice.Map(result.Users, transformUser),
		TotalCount: result.TotalCount,
	}, nil
}

func (q *AuthQuery) ListRoles(ctx context.Context) ([]Role, error) {
	roles, err := q.RBAC.GetAllRoles(ctx)
	if err != nil {
		return nil, err
	}

	return slice.Map(roles, transformRole), nil
}

func (q *AuthQuery) ListObjectActions() []AuthObjectAction {
	return slice.Map(q.RBAC.GetObjectActions(), func(objAct rbac.ObjectAction) AuthObjectAction {
		return AuthObjectAction{
			Namespace: objAct.Namespace,
			Object:    objAct.Object,
			Action:    objAct.Action,
		}
	})
}

type AuthMutation struct {
	User user.Service
	RBAC rbac.Repository
}

func (m *AuthMutation) SetUserRole(
	ctx context.Context,
	userID int,
	roleName string,
) (User, error) {
	user, err := m.User.SetRole(
		ctx,
		userID,
		roleName,
	)
	if err != nil {
		return User{}, err
	}

	return transformUser(user), nil
}

func (m *AuthMutation) SetUserEnabled(
	ctx context.Context,
	userID int,
	enabled bool,
) (User, error) {
	user, err := m.User.SetEnabled(
		ctx,
		userID,
		enabled,
	)
	if err != nil {
		return User{}, err
	}

	return transformUser(user), nil
}

func (m *AuthMutation) DeleteUser(ctx context.Context, userID int) (*string, error) {
	return nil, m.User.Delete(ctx, userID)
}

func (m *AuthMutation) PutRole(
	ctx context.Context,
	role string,
	objectActions []gen.AuthObjectActionInput,
) (*Role, error) {
	roleInfo, err := m.RBAC.PutRole(
		ctx,
		rbac.Role(role),
		slice.Map(objectActions, func(perm gen.AuthObjectActionInput) rbac.ObjectAction {
			return rbac.NewObjectAction(perm.Namespace, perm.Object, perm.Action)
		}),
	)
	if err != nil {
		return nil, err
	}

	gqlRole := transformRole(roleInfo)

	return &gqlRole, nil
}

func (m *AuthMutation) DeleteRole(
	ctx context.Context,
	role string,
) (*string, error) {
	return nil, m.RBAC.DeleteRole(
		ctx,
		rbac.Role(role),
	)
}

type Invitation struct {
	Code      string
	Role      string
	Email     *string
	CreatedBy *User
	ClaimedBy *User
	ExpiresAt *time.Time
	CreatedAt time.Time
}

type ListInvitationsResult struct {
	Invitations []Invitation
	TotalCount  int
}

func transformInvitation(inv model.Invitation) Invitation {
	var (
		email     *string
		createdBy *User
		claimedBy *User
		expiresAt *time.Time
	)

	if inv.Email.Valid {
		email = &inv.Email.String
	}

	if inv.CreatedByUser.ID > 0 {
		createdByV := transformUser(inv.CreatedByUser)
		createdBy = &createdByV
	}

	if inv.ClaimedByUser.ID > 0 {
		claimedByV := transformUser(inv.ClaimedByUser)
		claimedBy = &claimedByV
	}

	return Invitation{
		Code:      inv.Code,
		Role:      inv.RoleName,
		Email:     email,
		CreatedBy: createdBy,
		ClaimedBy: claimedBy,
		ExpiresAt: expiresAt,
		CreatedAt: inv.CreatedAt,
	}
}

func (m *AuthMutation) Invite(ctx context.Context, input gen.InviteInput) (Invitation, error) {
	usr, ok := auth.UserFromContext(ctx)
	if !ok {
		return Invitation{}, errors.New("no user found in context")
	}

	req := user.InviteRequest{
		CreatedBy: usr.ID,
	}

	if email, ok := input.Email.ValueOK(); ok {
		req.Email = *email
	}

	if expiration, ok := input.Expiration.ValueOK(); ok {
		req.Expiration = *expiration
	}

	inv, err := m.User.Invite(ctx, req)
	if err != nil {
		return Invitation{}, err
	}

	return transformInvitation(inv), nil
}

func (m *AuthMutation) DeleteInvitation(ctx context.Context, code string) (*string, error) {
	return nil, m.User.DeleteInvitation(ctx, code)
}

func (q *AuthQuery) ListInvitations(ctx context.Context, input *gen.ListInvitationsInput) (ListInvitationsResult, error) {
	var params user.ListInvitationsParams

	if input != nil {
		if pagination, ok := input.Pagination.ValueOK(); ok {
			if limit, ok := pagination.Limit.ValueOK(); ok {
				params.Limit = *limit
			}

			if offset, ok := pagination.Offset.ValueOK(); ok {
				params.Offset = *offset
			}

			if page, ok := pagination.Page.ValueOK(); ok {
				params.Page = *page
			}
		}
	}

	result, err := q.User.ListInvitations(ctx, params)
	if err != nil {
		return ListInvitationsResult{}, err
	}

	return ListInvitationsResult{
		TotalCount:  result.TotalCount,
		Invitations: slice.Map(result.Invitations, transformInvitation),
	}, nil
}
