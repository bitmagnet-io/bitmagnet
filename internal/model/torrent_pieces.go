package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (*TorrentPieces) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{
		UpdateAll: true,
	})
	return nil
}
