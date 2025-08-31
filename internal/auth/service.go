package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
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
)

type AuthService interface {
	Register(ctx context.Context, username, password string) (model.User, error)
	Login(ctx context.Context, username, password string) (LoginSuccess, error)
	GetUserByID(ctx context.Context, id int32) (model.User, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
}

type LoginSuccess struct {
	Token string
	User  model.User
}

type authService struct {
	JWTService
	database.DaoProvider
	PasswordPolicyService
}

func NewAuthService(
	daoProvider database.DaoProvider,
	jwtService JWTService,
	passwordPolicyService PasswordPolicyService,
) AuthService {
	return &authService{
		DaoProvider:           daoProvider,
		JWTService:            jwtService,
		PasswordPolicyService: passwordPolicyService,
	}
}

func (a *authService) Register(ctx context.Context, username, password string) (model.User, error) {
	dao, err := a.DaoProvider.Dao()
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrRegister, err)
	}

	// Validate password policy
	if err := a.PasswordPolicyService.ValidatePassword(password); err != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrRegister, fmt.Errorf("%w: %w", ErrPasswordPolicy, err))
	}

	// Check if user already exists
	existing, err := dao.WithContext(ctx).User.Where(dao.User.Username.Eq(username)).Count()
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrRegister, err)
	}

	if existing > 0 {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrRegister, ErrUserAlreadyExists)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrRegister, err)
	}

	// Create user
	user := &model.User{
		Username:  username,
		Password:  model.NewNullString(string(hashedPassword)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = dao.WithContext(ctx).User.Create(user)
	if err != nil {
		return model.User{}, fmt.Errorf("%w: %w: %w", Err, ErrRegister, err)
	}

	// Clear password from response
	user.Password = model.NullString{}

	return *user, nil
}

func (a *authService) Login(ctx context.Context, username, password string) (LoginSuccess, error) {
	dao, err := a.DaoProvider.Dao()
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

func (a *authService) GetUserByID(ctx context.Context, id int32) (model.User, error) {
	dao, err := a.DaoProvider.Dao()
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
	dao, err := a.DaoProvider.Dao()
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
