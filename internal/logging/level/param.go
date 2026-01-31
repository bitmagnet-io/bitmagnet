package level

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
)

var Param = param.MustNew(
	param.Dynamic(
		param.Description[Level]("The logging level"),
		param.EnumValues(LevelValues()...),
		param.Default(LevelInfo),
	),
)
