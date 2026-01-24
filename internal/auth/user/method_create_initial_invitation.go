package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gorm"
)

type InitialInvitationStatus int

const (
	InitialInvitationNotRequired InitialInvitationStatus = iota
	InitialInvitationCreated
	InitialInvitationUnclaimed
)

type InitialInvitation struct {
	model.Invitation
	Status InitialInvitationStatus
}

func (s *service) CreateInitialInvitation(ctx context.Context) (InitialInvitation, error) {
	var initialInvitation InitialInvitation

	err := s.DaoTransaction(func(tx *dao.Query) error {
		count, err := tx.WithContext(ctx).User.Where(
			tx.User.RoleName.Eq("admin"),
			tx.User.Enabled.Is(true),
		).Count()
		if err != nil {
			return err
		}

		if count > 0 {
			// we already have an admin user
			initialInvitation.Status = InitialInvitationNotRequired
			return nil
		}

		invitation, err := tx.WithContext(ctx).
			Invitation.
			Where(
				tx.Invitation.RoleName.Eq("admin"),
				tx.Invitation.CreatedBy.IsNull(),
				tx.Invitation.ClaimedBy.IsNull(),
				tx.Invitation.ExpiresAt.IsNull(),
			).
			First()
		if err == nil {
			initialInvitation.Invitation = *invitation
			initialInvitation.Status = InitialInvitationUnclaimed

			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		invitation = &model.Invitation{
			Code:     newInvitationCode(),
			RoleName: "admin",
		}

		err = tx.WithContext(ctx).Invitation.Create(invitation)
		if err != nil {
			return err
		}

		initialInvitation.Invitation = *invitation
		initialInvitation.Status = InitialInvitationCreated

		return nil
	})
	if err != nil {
		return InitialInvitation{}, fmt.Errorf("%w: %w: %w", Err, ErrInitialInvitation, err)
	}

	return initialInvitation, nil
}
