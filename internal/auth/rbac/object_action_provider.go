package rbac

import (
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type ObjectActionProvider func() []ObjectAction

func ObjectActionProviders(providers ...ObjectActionProvider) ObjectActionProvider {
	return func() []ObjectAction {
		objActs := slice.FlatMap(providers, func(provider ObjectActionProvider) []ObjectAction {
			return provider()
		})

		slices.SortFunc(objActs, func(a, b ObjectAction) int {
			return a.Compare(b)
		})

		return objActs
	}
}
