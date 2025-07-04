package manager

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
)

type manager struct {
	db database.GormDBProvider
}
