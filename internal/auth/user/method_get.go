package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gorm"
)

func (s *service) Get(ctx context.Context, id int) (model.User, error) {
	dao, err := s.Dao()
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrGet, err)
	}

	user, err := dao.WithContext(ctx).User.Where(dao.User.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrGet, ErrNotFound)
		}

		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrGet, err)
	}

	return *user, nil
}

func (s *service) GetByUsername(ctx context.Context, username string) (model.User, error) {
	dao, err := s.Dao()
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrGet, err)
	}

	user, err := dao.WithContext(ctx).User.Where(dao.User.Username.Eq(username)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrGet, ErrNotFound)
		}

		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrGet, err)
	}

	return *user, nil
}
