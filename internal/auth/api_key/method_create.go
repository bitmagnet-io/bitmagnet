package api_key

import (
	"context"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
)

type CreateRequest struct {
	UserID      int
	Name        string
	Permissions []rbac.ObjectAction
	Expiry      time.Duration
}

type CreateResult struct {
	ID        int
	APIKey    string
	Name      string
	ExpiresAt time.Time
}

func (s service) Create(ctx context.Context, req CreateRequest) (CreateResult, error) {
	secret := NewSecret()

	var expiresAt time.Time
	if req.Expiry > 0 {
		expiresAt = time.Now().Add(req.Expiry)
	}

	apiKeyID, err := s.repository.Create(ctx, req.UserID, req.Name, secret.Hash, req.Permissions, expiresAt)
	if err != nil {
		return CreateResult{}, fmt.Errorf("%w: %w: %w", Err, ErrCreate, err)
	}

	return CreateResult{
		ID: apiKeyID,
		APIKey: KeyData{
			ID:     apiKeyID,
			Secret: secret.Secret,
		}.Encode(),
		Name:      req.Name,
		ExpiresAt: expiresAt,
	}, nil
}
