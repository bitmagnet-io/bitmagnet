package dao

import (
	"fmt"
	"gorm.io/gorm/callbacks"
)

func (t torrentContent) CountEstimate() (int64, error) {
	db := t.UnderlyingDB()
	callbacks.BuildQuerySQL(db)
	query := db.Statement.SQL.String()
	args := db.Statement.Vars
	fmt.Printf("query: %s, args: %v", query, args)
	return 0, nil
}
