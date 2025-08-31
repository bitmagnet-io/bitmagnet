package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	Err                   = errors.New("auth")
	ErrRegister           = errors.New("registration failed")
	ErrLogin              = errors.New("login failed")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserHasNoPassword  = errors.New("user has no password")
	ErrIncorrectPassword  = errors.New("password is incorrect")
	ErrGenerateToken      = errors.New("generate token failed")
	ErrGetUser            = errors.New("get user failed")
	ErrPasswordPolicy     = errors.New("password policy validation failed")
	ErrUpdatePassword     = errors.New("update password failed")
	ErrTransaction        = errors.New("database transaction failed")
)

type AuthService interface {
	Register(ctx context.Context, username, password string) (model.User, error)
	Login(ctx context.Context, username, password string) (LoginSuccess, error)
	UpdatePassword(ctx context.Context, userID int32, currentPassword, newPassword string) error
	GetUserByID(ctx context.Context, id int32) (model.User, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
}

type LoginSuccess struct {
	Token string
	User  model.User
}

type authService struct {
	database.DaoTransactionProvider
	JWTService
	PasswordPolicyService
}

func NewAuthService(
	daoProvider database.DaoTransactionProvider,
	jwtService JWTService,
	passwordPolicyService PasswordPolicyService,
) AuthService {
	return &authService{
		DaoTransactionProvider: daoProvider,
		JWTService:             jwtService,
		PasswordPolicyService:  passwordPolicyService,
	}
}

func (a *authService) Register(ctx context.Context, username, password string) (model.User, error) {
	// Validate password policy
	if err := a.PasswordPolicyService.ValidatePassword(password); err != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrRegister, fmt.Errorf("%w: %w", ErrPasswordPolicy, err))
	}

	// Create user
	user := &model.User{
		Username:  username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Hash password
	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrRegister, err)
	} else {
		user.Password = model.NewNullString(string(hashedPassword))
	}

	var errUser error

	errTx := a.DaoTransaction(func(tx *dao.Query) error {
		// Check if user already exists
		if existing, err := tx.WithContext(ctx).User.Where(dao.User.Username.Eq(username)).Count(); err != nil {
			return err
		} else if existing > 0 {
			errUser = ErrUserAlreadyExists

			return nil
		}

		return tx.WithContext(ctx).User.Create(user)
	})

	if errTx != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w: %w", Err, ErrRegister, ErrTransaction, errTx)
	} else if errUser != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrRegister, errUser)
	}

	// Clear password from response
	user.Password = model.NullString{}

	return *user, nil
}

func (a *authService) Login(ctx context.Context, username, password string) (LoginSuccess, error) {
	dao, err := a.Dao()
	if err != nil {
		return LoginSuccess{}, fmt.Errorf("%w: %w: %w", Err, ErrLogin, err)
	}

	user, err := dao.WithContext(ctx).User.Where(dao.User.Username.Eq(username)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return LoginSuccess{}, fmt.Errorf("%w: %w: %w: %w", Err, ErrLogin, ErrInvalidCredentials, ErrUserNotFound)
		}

		return LoginSuccess{}, fmt.Errorf("%w: %w: %w", Err, ErrLogin, err)
	}

	if !user.Password.Valid {
		return LoginSuccess{}, fmt.Errorf("%w: %w: %w: %w", Err, ErrLogin, ErrInvalidCredentials, ErrUserHasNoPassword)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
	if err != nil {
		return LoginSuccess{}, fmt.Errorf("%w: %w: %w: %w", Err, ErrLogin, ErrInvalidCredentials, ErrIncorrectPassword)
	}

	token, err := a.JWTService.GenerateToken(user.ID, user.Username)
	if err != nil {
		return LoginSuccess{}, fmt.Errorf("%w: %w: %w: %w", Err, ErrLogin, ErrGenerateToken, err)
	}

	// Update last login
	dao.WithContext(ctx).User.Where(dao.User.ID.Eq(user.ID)).UpdateSimple(dao.User.LastLoginAt.Value(time.Now()))

	// Clear password from response
	user.Password = model.NullString{}

	return LoginSuccess{
		Token: token,
		User:  *user,
	}, nil
}

func (a *authService) UpdatePassword(ctx context.Context, userID int32, currentPassword, newPassword string) error {
	// Validate new password policy
	if err := a.PasswordPolicyService.ValidatePassword(newPassword); err != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrUpdatePassword, fmt.Errorf("%w: %w", ErrPasswordPolicy, err))
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrUpdatePassword, err)
	}

	var errUser error

	errTx := a.DaoTransaction(func(tx *dao.Query) error {
		// Get the user
		user, err := tx.WithContext(ctx).User.Where(dao.User.ID.Eq(userID)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				errUser = ErrUserNotFound

				return nil
			}

			return err
		}

		// Check if user has a password
		if user.Password.Valid {
			// Verify current password
			err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(currentPassword))
			if err != nil {
				errUser = fmt.Errorf("%w: %w", ErrIncorrectPassword, err)

				return nil
			}
		} else if currentPassword != "" {
			errUser = ErrIncorrectPassword

			return nil
		}

		// Update password
		_, err = tx.WithContext(ctx).User.Where(dao.User.ID.Eq(userID)).UpdateSimple(
			dao.User.Password.Value(model.NewNullString(string(hashedPassword))),
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

func (a *authService) GetUserByID(ctx context.Context, id int32) (model.User, error) {
	dao, err := a.Dao()
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrGetUser, err)
	}

	user, err := dao.WithContext(ctx).User.Where(dao.User.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrGetUser, ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrGetUser, err)
	}

	// Clear password from response
	user.Password = model.NullString{}

	return *user, nil
}

func (a *authService) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	dao, err := a.Dao()
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrGetUser, err)
	}

	user, err := dao.WithContext(ctx).User.Where(dao.User.Username.Eq(username)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrGetUser, ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrGetUser, err)
	}

	// Clear password from response
	user.Password = model.NullString{}

	return *user, nil
}
