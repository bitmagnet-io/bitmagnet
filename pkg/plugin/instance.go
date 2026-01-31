package plugin

import "github.com/bitmagnet-io/bitmagnet/internal/ref"

type Instance interface {
	Ref() ref.Ref
	Enabled() bool
	DependsOn() []ref.Ref
	RequiredBy() []ref.Ref
	LocalizedContentProvider
}
