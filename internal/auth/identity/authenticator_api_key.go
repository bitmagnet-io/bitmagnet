package identity

import (
	"context"
	"errors"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/api_key"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
)

type authenticatorAPIKey struct {
	apiKeyService api_key.Service
	rbac          rbac.Service
}

func (a authenticatorAPIKey) Authenticate(ctx context.Context, token string) (Identity, bool, error) {
	apiKey, err := a.apiKeyService.Auth(ctx, token)
	if err != nil {
		return nil, !errors.Is(err, api_key.ErrDecode), err
	}

	anon, err := a.rbac.GetRole(ctx, rbac.RoleAnon)
	if err != nil {
		return nil, true, err
	}

	return APIKey{
		APIKey:   apiKey,
		anon:     anon,
		enforcer: a.rbac,
	}, true, nil
}
