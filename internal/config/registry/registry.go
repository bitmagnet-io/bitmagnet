package registry

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
)

type (
	Param struct {
		Ref    ref.Ref
		Plugin ref.Ref
		param.Untyped
	}
)
