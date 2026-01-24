package auth

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
)

var ErrUnauthorized = errors.New("unauthorized")

type unauthorizedError struct {
	objAct rbac.ObjectAction
}

func (unauthorizedError) Error() string {
	return ErrUnauthorized.Error()
}

func (unauthorizedError) Unwrap() error {
	return ErrUnauthorized
}

func (e unauthorizedError) GraphQLExtensions() map[string]any {
	return map[string]any{
		"namespace": Namespace,
		"object":    e.objAct.Object,
		"action":    e.objAct.Action,
	}
}

type Directive func(
	ctx context.Context,
	obj any,
	next graphql.Resolver,
	object string,
	action string,
) (res any, err error)

const Namespace = "graphql"

func NewDirective() Directive {
	return func(
		ctx context.Context,
		_ any,
		next graphql.Resolver,
		object string,
		action string,
	) (res any, err error) {
		allow := false

		objAct := rbac.ObjectAction{
			Namespace: Namespace,
			Object:    object,
			Action:    action,
		}

		identity, ok := IdentityFromContext(ctx)
		if ok {
			var err error

			allow, err = identity.Enforce(ctx, objAct)
			if err != nil {
				return nil, err
			}
		}

		if !allow {
			return nil, unauthorizedError{
				objAct: objAct,
			}
		}

		return next(ctx)
	}
}
