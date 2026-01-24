package plugin

import "github.com/bitmagnet-io/bitmagnet/internal/ref"

type Infos []Info

type Info struct {
	Ref        ref.Ref
	Enabled    bool
	DependsOn  []ref.Ref
	RequiredBy []ref.Ref
}
