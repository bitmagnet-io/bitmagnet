package auth

import (
	"maps"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/directive"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

func Permissions() []rbac.Permission {
	anon := rbac.SubjectRole{
		Role: rbac.RoleAnon,
	}
	user := rbac.SubjectRole{
		Role: rbac.RoleUser,
	}

	return append(
		[]rbac.Permission{
			// User-specific permissions:
			rbac.NewPermission(user, rbac.NewObjectAction(Namespace, "torrent", "query")),
		},
		slice.FlatMap([]rbac.SubjectRole{
			anon,
			user,
		}, func(role rbac.SubjectRole) []rbac.Permission {
			// Permissions granted to all users and anons:
			return []rbac.Permission{
				rbac.NewPermission(role, rbac.NewObjectAction(Namespace, "self", "mutate")),
				rbac.NewPermission(role, rbac.NewObjectAction(Namespace, "self", "query")),
				rbac.NewPermission(role, rbac.NewObjectAction(Namespace, "health", "query")),
				rbac.NewPermission(role, rbac.NewObjectAction(Namespace, "version", "query")),
			}
		})...)
}

func AuthObjectActions(directives directive.AuthDirectives) []rbac.ObjectAction {
	return slice.Map(slices.Collect(maps.Keys(directives)), func(dir directive.AuthDirective) rbac.ObjectAction {
		return rbac.NewObjectAction("graphql", dir.Object, dir.Action)
	})
}
