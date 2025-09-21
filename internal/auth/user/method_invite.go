package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type InviteRequest struct {
	Email     string
	Role      string
	CreatedBy int
	Expiry    time.Duration
}

func (s *service) Invite(ctx context.Context, request InviteRequest) (model.Invitation, error) {
	invitation := model.Invitation{
		RoleName: request.Role,
	}

	if request.Email != "" {
		if !regexEmail.MatchString(request.Email) {
			return model.Invitation{}, fmt.Errorf("%w: %w: %w", Err, ErrInvite, ErrEmailInvalid)
		}

		invitation.Email = model.NewNullString(request.Email)
	}

	if request.CreatedBy > 0 {
		invitation.CreatedBy = model.NewNullInt(request.CreatedBy)
	}

	if request.Expiry > 0 {
		invitation.ExpiresAt = sql.NullTime{
			Time:  time.Now().Add(request.Expiry),
			Valid: true,
		}
	}

	invitation.Code = newInvitationCode()

	err := s.DaoTransaction(func(tx *dao.Query) error {
		err := tx.WithContext(ctx).Invitation.Create(&invitation)
		if err != nil {
			return err
		}

		if request.CreatedBy > 0 {
			createdBy, err := tx.WithContext(ctx).User.Where(tx.User.ID.Eq(request.CreatedBy)).First()
			if err != nil {
				return err
			}

			invitation.CreatedByUser = *createdBy
		}

		return nil
	})

	if err != nil {
		return model.Invitation{}, fmt.Errorf("%w: %w: %w: %w", Err, ErrInvite, ErrTransaction, err)
	}

	return invitation, nil
}
