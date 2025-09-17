package auth

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

var ErrUnauthorized = errors.New("unauthorized")

type unauthorizedError struct {
	subjects []rbac.Subject
	objAct   rbac.ObjectAction
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
		"subjects": slice.Map(e.subjects, func(subj rbac.Subject) map[string]string {
			return map[string]string{
				"type": string(subj.SubjectType()),
				"name": subj.SubjectName(),
			}
		}),
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

func NewDirective(enforcer rbac.Enforcer) Directive {
	return func(
		ctx context.Context,
		_ any,
		next graphql.Resolver,
		object string,
		action string,
	) (res any, err error) {
		subjects := subjectsFromContext(ctx)

		objAct := rbac.NewObjectAction(Namespace, object, action)

		allow, err := enforcer.EnforceAny(
			ctx,
			subjects,
			objAct,
		)
		if err != nil {
			return nil, err
		}

		if !allow {
			return nil, unauthorizedError{
				subjects: subjects,
				objAct:   objAct,
			}
			// return nil, &gqlerror.Error{
			// 	Message: "Unauthorized",
			// 	Extensions: map[string]any{
			// 		"code":      "UNAUTHORIZED",
			// 		"namespace": Namespace,
			// 		"object":    object,
			// 		"action":    action,
			// 		"subjects": slice.Map(subjects, func(subj rbac.Subject) map[string]string {
			// 			return map[string]string{
			// 				"type": string(subj.SubjectType()),
			// 				"name": subj.SubjectName(),
			// 			}
			// 		}),
			// 	},
			// }
		}

		return next(ctx)
	}
}

func subjectsFromContext(ctx context.Context) []rbac.Subject {
	_, roles := UserRolesFromContext(ctx)

	return slice.Map(roles, func(role rbac.Role) rbac.Subject {
		return rbac.SubjectRole{Role: role}
	})
}
