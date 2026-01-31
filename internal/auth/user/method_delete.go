package user

import (
	"context"
	"fmt"
)

func (s *service) Delete(ctx context.Context, userID int) error {
	dao, err := s.Dao()
	if err != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrDelete, err)
	}

	info, err := dao.WithContext(ctx).User.Where(dao.User.ID.Eq(userID)).Delete()
	if err != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrDelete, err)
	}

	if info.RowsAffected == 0 {
		return fmt.Errorf("%w: %w: %w", Err, ErrDelete, ErrNotFound)
	}

	return nil
}
