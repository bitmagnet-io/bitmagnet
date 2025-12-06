package api_key

import (
	"context"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (s service) Auth(ctx context.Context, key string) (model.APIKey, error) {
	keyData := &KeyData{}

	if err := keyData.Decode(key); err != nil {
		return model.APIKey{}, fmt.Errorf("%w: %w: %w: %w", Err, ErrAuth, ErrDecode, err)
	}

	apiKey, err := s.repository.Get(ctx, keyData.ID)

	if err == nil && bcrypt.CompareHashAndPassword(apiKey.Hash, keyData.Secret) != nil {
		err = ErrMismatch
	}

	if err != nil {
		return model.APIKey{}, fmt.Errorf("%w: %w: %w", Err, ErrAuth, err)
	}

	if apiKey.ExpiresAt.Valid && apiKey.ExpiresAt.Time.Before(time.Now()) {
		return model.APIKey{}, fmt.Errorf("%w: %w: %w", Err, ErrAuth, ErrExpired)
	}

	apiKey.Hash = nil
	apiKey.User.Password = nil

	return apiKey, nil
}
