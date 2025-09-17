package user

import (
	"context"
	"fmt"
)

func (s *service) DeleteInvitation(ctx context.Context, code string) error {
	dao, err := s.Dao()
	if err != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrDelete, err)
	}

	info, err := dao.WithContext(ctx).Invitation.Where(dao.Invitation.Code.Eq(code)).Delete()
	if err != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrDeleteInvitation, err)
	}

	if info.RowsAffected == 0 {
		return fmt.Errorf("%w: %w: %w", Err, ErrDeleteInvitation, ErrInvitationNotFound)
	}

	return nil
}
