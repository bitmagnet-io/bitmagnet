// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameTorrentSource = "torrent_sources"

// TorrentSource mapped from table <torrent_sources>
type TorrentSource struct {
	Key       string    `gorm:"column:key;primaryKey;<-:create" json:"key"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;not null;<-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updatedAt"`
}

// TableName TorrentSource's table name
func (*TorrentSource) TableName() string {
	return TableNameTorrentSource
}
