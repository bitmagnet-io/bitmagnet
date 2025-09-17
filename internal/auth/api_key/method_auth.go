package api_key

import (
	"context"
	"errors"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *service) Auth(ctx context.Context, key string) (model.APIKey, error) {
	keyData := &keyData{}

	if err := keyData.decode(key); err != nil {
		return model.APIKey{}, fmt.Errorf("%w: %w: %w: %w", Err, ErrAuth, ErrDecode, err)
	}

	dao, err := s.dao.Dao()
	if err != nil {
		return model.APIKey{}, fmt.Errorf("%w: %w: %w", Err, ErrAuth, err)
	}

	apiKey, err := dao.APIKey.
		WithContext(ctx).
		Where(
			dao.APIKey.ID.Eq(keyData.id),
		).
		Preload(
			dao.APIKey.User,
			dao.APIKey.Permissions,
		).
		First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrNotFound
	}

	if err == nil && bcrypt.CompareHashAndPassword(apiKey.Hash, keyData.secret) != nil {
		err = ErrMismatch
	}

	if err != nil {
		return model.APIKey{}, fmt.Errorf("%w: %w: %w", Err, ErrAuth, ErrNotFound)
	}

	apiKey.Hash = nil
	apiKey.User.Password = nil

	return *apiKey, nil
}
