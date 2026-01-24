package gqlmodel

import (
	"context"
	"errors"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/api_key"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/identity"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/user"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/auth"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type SelfQuery struct {
	RBAC   rbac.Repository
	User   user.Service
	APIKey api_key.Service
}

type Self struct {
	User        *User
	APIKey      *APIKey
	Permissions []AuthObjectAction
}

type APIKey struct {
	ID        int
	Name      string
	UserID    int
	User      User
	ExpiresAt *time.Time
	CreatedAt time.Time
}

func transformAPIKey(apiKey model.APIKey) APIKey {
	var expiresAt *time.Time
	if apiKey.ExpiresAt.Valid {
		expiresAt = &apiKey.ExpiresAt.Time
	}

	return APIKey{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		UserID:    apiKey.UserID,
		User:      transformUser(apiKey.User),
		ExpiresAt: expiresAt,
		CreatedAt: apiKey.CreatedAt,
	}
}

func transformSelf(self identity.Self) Self {
	var user *User

	if self.User != nil {
		transformed := transformUser(*self.User)
		user = &transformed
	}

	var apiKey *APIKey

	if self.APIKey != nil {
		transformed := transformAPIKey(*self.APIKey)
		apiKey = &transformed
	}

	return Self{
		User:        user,
		APIKey:      apiKey,
		Permissions: self.Permissions,
	}
}

func (*SelfQuery) Identity(ctx context.Context) (Self, error) {
	identity, ok := auth.IdentityFromContext(ctx)
	if !ok {
		return Self{}, errors.New("no identity in context")
	}

	return transformSelf(identity.Self()), nil
}

func (q *SelfQuery) APIKeys(ctx context.Context) ([]APIKey, error) {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, errors.New("no user in context")
	}

	result, err := q.APIKey.List(ctx, api_key.ListRequest{
		UserID: user.ID,
	})
	if err != nil {
		return nil, err
	}

	return slice.Map(result.APIKeys, transformAPIKey), nil
}

func (q *SelfQuery) PasswordEntropy(password string) user.PasswordEntropyResult {
	return q.User.PasswordEntropy(password)
}

type SelfMutation struct {
	User   user.Service
	RBAC   rbac.Repository
	APIKey api_key.Service
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

func (m *SelfMutation) CreateAPIKey(ctx context.Context, input gen.CreateAPIKeyInput) (gen.CreateAPIKeyResult, error) {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return gen.CreateAPIKeyResult{}, errors.New("no user in context")
	}

	var expiry time.Duration

	if inputExpiry, ok := input.Expiry.ValueOK(); ok {
		expiry = *inputExpiry
	}

	result, err := m.APIKey.Create(ctx, api_key.CreateRequest{
		UserID: user.ID,
		Name:   input.Name,
		Permissions: slice.Map(input.Permissions, func(perm gen.AuthObjectActionInput) rbac.ObjectAction {
			return rbac.ObjectAction{
				Namespace: perm.Namespace,
				Object:    perm.Object,
				Action:    perm.Action,
			}
		}),
		Expiry: expiry,
	})
	if err != nil {
		return gen.CreateAPIKeyResult{}, err
	}

	var expiresAt *time.Time
	if !result.ExpiresAt.IsZero() {
		expiresAt = &result.ExpiresAt
	}

	return gen.CreateAPIKeyResult{
		ID:        result.ID,
		APIKey:    result.APIKey,
		Name:      result.Name,
		ExpiresAt: expiresAt,
	}, nil
}

func (m *SelfMutation) DeleteAPIKey(ctx context.Context, id int) (*string, error) {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, errors.New("no user in context")
	}

	return nil, m.APIKey.Delete(ctx, api_key.DeleteRequest{
		UserID:   user.ID,
		APIKeyID: id,
	})
}
