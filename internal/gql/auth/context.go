package auth

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/http_auth"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	internal_httpserver "github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func UserRolesFromContext(ctx context.Context) (*model.User, []rbac.Role) {
	user, ok := UserFromContext(ctx)
	if ok {
		return &user, []rbac.Role{rbac.Role(user.RoleName)}
	}

	return nil, []rbac.Role{
		rbac.RoleAnon,
	}
}

func UserFromContext(ctx context.Context) (model.User, bool) {
	ginCtx, ok := internal_httpserver.GinContextFromContext(ctx)
	if !ok {
		return model.User{}, false
	}

	return http_auth.GetCurrentUser(ginCtx)
}
