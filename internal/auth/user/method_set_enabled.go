package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gorm"
)

func (s *service) SetEnabled(ctx context.Context, userID int, enabled bool) (model.User, error) {
	var (
		user    *model.User
		userErr error
	)

	err := s.DaoTransaction(func(tx *dao.Query) error {
		_, err := tx.WithContext(ctx).
			User.
			Where(tx.User.ID.Eq(userID)).
			UpdateColumn(tx.User.Enabled, enabled)
		if err != nil {
			return err
		}

		user, err = tx.WithContext(ctx).User.Where(tx.User.ID.Eq(userID)).First()

		if errors.Is(err, gorm.ErrRecordNotFound) {
			userErr = ErrNotFound
			return nil
		}

		return err
	})

	if err == nil {
		err = userErr
	}

	if err != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrSetEnabled, err)
	}

	return *user, nil
}
