package api_key

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"golang.org/x/crypto/bcrypt"
)

type CreateRequest struct {
	UserID      int
	Name        string
	Permissions []rbac.ObjectAction
	Expiry      time.Duration
}

type CreateResult struct {
	APIKeyID  int
	APIKey    string
	ExpiresAt time.Time
}

func (s *service) Create(ctx context.Context, req CreateRequest) (CreateResult, error) {
	secret := newSecret()

	record := model.APIKey{
		UserID: req.UserID,
		Hash:   secret.hash,
		Permissions: slice.Map(req.Permissions, func(perm rbac.ObjectAction) model.APIKeyPermission {
			return model.APIKeyPermission{
				Namespace: perm.Namespace,
				Object:    perm.Object,
				Action:    perm.Action,
			}
		}),
	}

	if req.Expiry > 0 {
		record.ExpiresAt = sql.NullTime{Time: time.Now().Add(req.Expiry), Valid: true}
	}

	dao, err := s.dao.Dao()
	if err != nil {
		return CreateResult{}, fmt.Errorf("%w: %w: %w", Err, ErrCreate, err)
	}

	err = dao.APIKey.WithContext(ctx).Create(&record)
	if err != nil {
		return CreateResult{}, fmt.Errorf("%w: %w: %w", Err, ErrCreate, err)
	}

	result := CreateResult{
		APIKeyID: record.ID,
		APIKey: keyData{
			id:     record.ID,
			secret: secret.secret,
		}.encode(),
	}

	if record.ExpiresAt.Valid {
		result.ExpiresAt = record.ExpiresAt.Time
	}

	return result, nil
}

type secret struct {
	secret []byte
	hash   []byte
}

func newSecret() secret {
	bytes := make([]byte, secretLength)
	_, _ = rand.Read(bytes)

	hash, _ := bcrypt.GenerateFromPassword(
		bytes,
		bcrypt.DefaultCost,
	)

	return secret{
		secret: bytes,
		hash:   hash,
	}
}
