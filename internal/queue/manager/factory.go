package manager

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
)

func New(
	db database.GormDBProvider,
) Manager {
	return manager{
		db: db,
	}
}
