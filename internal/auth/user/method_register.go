package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username       string
	Password       string
	Email          string
	InvitationCode string
}

func (s *service) Register(ctx context.Context, request RegisterRequest) (model.User, error) {
	if !regexUsername.MatchString(request.Username) {
		return model.User{}, fmt.Errorf(
			"%w: %w: %w: must match %s",
			Err,
			ErrRegister,
			ErrUsernameInvalid,
			regexUsername.String(),
		)
	}

	user := model.User{
		Username:  request.Username,
		RoleName:  rbac.RoleUser.String(),
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if request.InvitationCode == "" && s.invitationRequired.Get() {
		return model.User{}, ErrInvitationCodeMissing
	}

	if request.Email == "" && s.emailRequired.Get() {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrRegister, ErrEmailMissing)
	}

	if request.Email != "" {
		if !regexEmail.MatchString(request.Email) {
			return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrRegister, ErrEmailInvalid)
		}

		user.Email = model.NewNullString(request.Email)
	}

	if !s.PasswordEntropy(request.Password).Valid {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrRegister, ErrPasswordInsufficientEntropy)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		int(s.passwordHashingCost.Get()),
	)
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrRegister, err)
	}

	user.Password = hashedPassword

	var errUser error

	errTx := s.DaoTransaction(func(tx *dao.Query) error {
		// Check invitation validity
		if request.InvitationCode != "" {
			invitation, err := tx.WithContext(ctx).Invitation.
				Where(tx.Invitation.Code.Eq(request.InvitationCode)).
				First()
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					errUser = ErrInvitationNotFound

					return nil
				}

				return err
			}

			if invitation.ClaimedBy.Valid {
				errUser = ErrInvitationClaimed
				return nil
			}

			if invitation.ExpiresAt.Valid && invitation.ExpiresAt.Time.After(time.Now()) {
				errUser = ErrInvitationExpired
				return nil
			}

			user.RoleName = invitation.RoleName
		}

		// Check if user already exists
		if existing, err := tx.WithContext(ctx).
			User.Where(tx.User.Username.Eq(request.Username)).
			Count(); err != nil {
			return err
		} else if existing > 0 {
			errUser = ErrAlreadyExists

			return nil
		}

		err := tx.WithContext(ctx).User.Create(&user)
		if err != nil {
			return err
		}

		if request.InvitationCode != "" {
			_, err = tx.WithContext(ctx).
				Invitation.
				Where(tx.Invitation.Code.Eq(request.InvitationCode)).
				UpdateColumn(tx.Invitation.ClaimedBy, user.ID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if errTx != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w: %w", Err, ErrRegister, ErrTransaction, errTx)
	} else if errUser != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrRegister, errUser)
	}

	// Clear password from response
	user.Password = nil

	return user, nil
}
