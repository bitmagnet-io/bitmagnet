package auth

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/http_auth"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/identity"
	internal_httpserver "github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func IdentityFromContext(ctx context.Context) (identity.Identity, bool) {
	ginCtx, ok := internal_httpserver.GinContextFromContext(ctx)
	if !ok {
		return nil, false
	}

	return http_auth.GetIdentity(ginCtx)
}

func UserFromContext(ctx context.Context) (model.User, bool) {
	identity, ok := IdentityFromContext(ctx)
	if !ok {
		return model.User{}, false
	}

	self := identity.Self()

	if self.User == nil {
		return model.User{}, false
	}

	return *self.User, true
}
