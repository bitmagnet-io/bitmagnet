package level

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

var Param = param.MustNew(
	param.WithDynamic(
		param.WithEnumValues(LevelValues()...),
		param.WithDefault(LevelInfo),
	),
)
