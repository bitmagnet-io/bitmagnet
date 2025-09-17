package user

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func (s *service) SetRole(ctx context.Context, userID int, roleName string) (model.User, error) {
	var (
		user    *model.User
		userErr error
	)

	err := s.DaoTransaction(func(tx *dao.Query) error {
		// first check the user and role exist
		count, err := tx.WithContext(ctx).User.Where(tx.User.ID.Eq(userID)).Count()
		if err != nil {
			return err
		}

		if count == 0 {
			userErr = ErrNotFound
			return nil
		}

		count, err = tx.WithContext(ctx).Role.Where(tx.Role.Name.Eq(roleName)).Count()
		if err != nil {
			return err
		}

		if count == 0 {
			userErr = ErrRoleNotFound
			return nil
		}

		_, err = tx.User.
			WithContext(ctx).
			UpdateColumn(tx.User.RoleName, roleName)
		if err != nil {
			return err
		}

		user, err = tx.WithContext(ctx).
			User.
			Where(tx.User.ID.Eq(userID)).
			First()

		return err
	})

	if err == nil {
		err = userErr
	}

	if err != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrAssignRole, err)
	}

	return *user, nil
}
