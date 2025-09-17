package api_key

import (
	"context"
	"fmt"
)

func (s *service) Delete(ctx context.Context, id int) error {
	dao, err := s.dao.Dao()
	if err != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrDelete, err)
	}

	info, err := dao.APIKey.WithContext(ctx).Where(dao.APIKey.ID.Eq(id)).Delete()
	if err != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrDelete, err)
	}

	if info.RowsAffected < 1 {
		return fmt.Errorf("%w: %w: %w", Err, ErrDelete, ErrNotFound)
	}

	return nil
}
