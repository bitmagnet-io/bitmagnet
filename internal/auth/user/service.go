package user

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/jwt"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"golang.org/x/time/rate"
)

type Service interface {
	// CreateInitialInvitation checks if there is any admin user or admin invitation.
	// If none is found, it creates an admin invitation.
	CreateInitialInvitation(ctx context.Context) (InitialInvitation, error)
	Invite(ctx context.Context, request InviteRequest) (model.Invitation, error)
	Register(ctx context.Context, request RegisterRequest) (model.User, error)
	Login(ctx context.Context, username, password string) (LoginResult, error)
	SetRole(ctx context.Context, userID int, roleName string) (model.User, error)
	UpdatePassword(ctx context.Context, userID int, currentPassword, newPassword string) error
	Get(ctx context.Context, userID int) (model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
	List(ctx context.Context, params ListUsersParams) (ListUsersResult, error)
	Delete(ctx context.Context, userID int) error
	SetEnabled(ctx context.Context, userID int, enabled bool) (model.User, error)
	ListInvitations(ctx context.Context, params ListInvitationsParams) (ListInvitationsResult, error)
	DeleteInvitation(ctx context.Context, code string) error
	PasswordEntropy(password string) PasswordEntropyResult
}

type service struct {
	database.DaoTransactionProvider
	jwtService          jwt.Service
	invitationRequired  *atomic.Value[InvitationRequired]
	emailRequired       *atomic.Value[EmailRequired]
	emailVerification   *atomic.Value[EmailVerification]
	passwordMinEntropy  *atomic.Value[PasswordMinEntropy]
	passwordHashingCost *atomic.Value[PasswordHashingCost]
	loginLimiter        *rate.Limiter
}

func NewService(
	daoProvider database.DaoTransactionProvider,
	jwtService jwt.Service,
	invitationRequired *atomic.Value[InvitationRequired],
	emailRequired *atomic.Value[EmailRequired],
	emailVerification *atomic.Value[EmailVerification],
	passwordMinEntropy *atomic.Value[PasswordMinEntropy],
	passwordHashingCost *atomic.Value[PasswordHashingCost],
	loginRequestsPerMinute *atomic.Value[LoginRequestsPerMinute],
	loginRequestBurst *atomic.Value[LoginRequestBurst],
) Service {
	return &service{
		DaoTransactionProvider: daoProvider,
		jwtService:             jwtService,
		invitationRequired:     invitationRequired,
		emailRequired:          emailRequired,
		emailVerification:      emailVerification,
		passwordMinEntropy:     passwordMinEntropy,
		passwordHashingCost:    passwordHashingCost,
		loginLimiter: newLoginLimiter(
			loginRequestsPerMinute,
			loginRequestBurst,
		),
	}
}
