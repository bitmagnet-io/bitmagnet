package dao

import (
	"gorm.io/gorm/callbacks"
)

func (t torrentContent) CountEstimate() (int64, error) {
	db := t.UnderlyingDB()
	callbacks.BuildQuerySQL(db)
	return 0, nil
}
