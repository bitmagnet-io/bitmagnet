package identity

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type Identity interface {
	Self() Self
	Enforce(context.Context, rbac.ObjectAction) (bool, error)
}

type Self struct {
	User        *model.User
	APIKey      *model.APIKey
	Permissions []rbac.ObjectAction
}
