package gqlmodel

import (
	"context"
	"errors"

	"github.com/bitmagnet-io/bitmagnet/internal/auth"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
)

type AuthMutation struct {
	auth.AuthService
}

func (m *AuthMutation) Register(ctx context.Context, username, password string) (gen.AuthRegisterResult, error) {
	user, err := m.AuthService.Register(ctx, username, password)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			return gen.AuthRegisterFailure{
				Error: gen.AuthRegisterErrorTypeUserAlreadyExists,
			}, nil
		}
		if errors.Is(err, auth.ErrPasswordPolicy) {
			return gen.AuthRegisterFailure{
				Error: gen.AuthRegisterErrorTypePasswordPolicyViolation,
			}, nil
		}

		return nil, err
	}

	return gen.AuthRegisterSuccess{
		User: user,
	}, nil
}

func (m *AuthMutation) Login(ctx context.Context, username, password string) (gen.AuthLoginResult, error) {
	success, err := m.AuthService.Login(ctx, username, password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return gen.AuthLoginFailure{
				Error: gen.AuthLoginErrorTypeInvalidCredentials,
			}, nil
		}

		return nil, err
	}

	return gen.AuthLoginSuccess{
		Token: success.Token,
		User:  success.User,
	}, nil
}
