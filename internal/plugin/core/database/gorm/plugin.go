package gorm

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database"
)

var Ref = database.Ref.MustSub("gorm")
