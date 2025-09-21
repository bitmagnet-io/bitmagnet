package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

type LoginResult struct {
	Token string
	User  model.User
}

func (s *service) Login(ctx context.Context, username, password string) (LoginResult, error) {
	if err := s.loginLimiter.Wait(ctx); err != nil {
		return LoginResult{}, fmt.Errorf("%w: %w", Err, ErrLoginRequestLimiter)
	}

	dao, err := s.Dao()
	if err != nil {
		return LoginResult{}, fmt.Errorf("%w: %w: %w", Err, ErrLogin, err)
	}

	user, err := dao.WithContext(ctx).
		User.
		Where(dao.User.Username.Eq(username)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return LoginResult{}, fmt.Errorf("%w: %w: %w: %w", Err, ErrLogin, ErrCredentialsInvalid, ErrNotFound)
		}

		return LoginResult{}, fmt.Errorf("%w: %w: %w", Err, ErrLogin, err)
	}

	if len(user.Password) == 0 {
		return LoginResult{}, fmt.Errorf("%w: %w: %w: %w", Err, ErrLogin, ErrCredentialsInvalid, ErrHasNoPassword)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return LoginResult{}, fmt.Errorf("%w: %w: %w: %w", Err, ErrLogin, ErrCredentialsInvalid, ErrIncorrectPassword)
	}

	if !user.Enabled {
		return LoginResult{}, fmt.Errorf("%w: %w: %w", Err, ErrLogin, ErrDisabled)
	}

	token, err := s.jwtService.Generate(user.ID, user.Username)
	if err != nil {
		return LoginResult{}, fmt.Errorf("%w: %w: %w: %w", Err, ErrLogin, ErrGenerateToken, err)
	}

	// Update last login
	_, _ = dao.WithContext(ctx).
		User.
		Where(dao.User.ID.Eq(user.ID)).
		UpdateSimple(
			dao.User.LastLoginAt.Value(sql.NullTime{Time: time.Now(), Valid: true}),
		)

	// Clear password from response
	user.Password = nil

	return LoginResult{
		Token: token,
		User:  *user,
	}, nil
}

func (rpm LoginRequestsPerMinute) limit() rate.Limit {
	return rate.Every(time.Minute / time.Duration(rpm))
}

func newLoginLimiter(
	rpm *atomic.Value[LoginRequestsPerMinute],
	burst *atomic.Value[LoginRequestBurst],
) *rate.Limiter {
	currentRpm, currentBurst := rpm.Get(), burst.Get()
	limiter := rate.NewLimiter(
		currentRpm.limit(),
		int(currentBurst),
	)

	rpm.Subscribe(func(rpm LoginRequestsPerMinute) {
		if rpm != currentRpm {
			currentRpm = rpm
			limiter.SetLimit(currentRpm.limit())
		}
	})

	burst.Subscribe(func(burst LoginRequestBurst) {
		if burst != currentBurst {
			currentBurst = burst
			limiter.SetBurst(int(currentBurst))
		}
	})

	return limiter
}
