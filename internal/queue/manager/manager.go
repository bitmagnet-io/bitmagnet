package manager

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"gorm.io/gorm"
)

type manager struct {
	dao *dao.Query
	db  *gorm.DB
}
