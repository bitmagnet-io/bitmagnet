package auth

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/user"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"go.uber.org/zap"
)

type initialInvitationWorker struct {
	userService user.Service
	logger      *zap.Logger
}

func (i *initialInvitationWorker) Runner() runner.Runner {
	return runner.SimpleRunner(func(ctx context.Context) error {
		result, err := i.userService.CreateInitialInvitation(ctx)
		if err != nil {
			return err
		}

		switch result.Status {
		case user.InitialInvitationCreated:
			i.logger.Warn("initial invitation created", zap.String("code", result.Code))
		case user.InitialInvitationUnclaimed:
			i.logger.Warn("initial invitation unclaimed", zap.String("code", result.Code))
		default:
		}

		return nil
	})
}
