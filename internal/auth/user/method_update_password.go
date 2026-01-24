package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *service) UpdatePassword(ctx context.Context, userID int, currentPassword, newPassword string) error {
	if !s.PasswordEntropy(newPassword).Valid {
		return fmt.Errorf("%w: %w: %w", Err, ErrUpdatePassword, ErrPasswordInsufficientEntropy)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrUpdatePassword, err)
	}

	var errUser error

	errTx := s.DaoTransaction(func(tx *dao.Query) error {
		user, err := tx.WithContext(ctx).User.Where(tx.User.ID.Eq(userID)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				errUser = ErrNotFound

				return nil
			}

			return err
		}

		if len(user.Password) > 0 {
			err = bcrypt.CompareHashAndPassword(user.Password, []byte(currentPassword))
			if err != nil {
				errUser = fmt.Errorf("%w: %w", ErrIncorrectPassword, err)

				return nil
			}
		} else if currentPassword != "" {
			errUser = ErrIncorrectPassword

			return nil
		}

		_, err = tx.WithContext(ctx).User.Where(tx.User.ID.Eq(userID)).UpdateSimple(
			tx.User.Password.Value(hashedPassword),
		)

		return err
	})
	if errTx != nil {
		return fmt.Errorf("%w: %w: %w: %w", Err, ErrUpdatePassword, ErrTransaction, errTx)
	} else if errUser != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrUpdatePassword, errUser)
	}

	return nil
}
