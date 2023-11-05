// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

const TableNameTorrentFile = "torrent_files"

// TorrentFile mapped from table <torrent_files>
type TorrentFile struct {
	InfoHash  protocol.ID `gorm:"column:info_hash;primaryKey;<-:create" json:"infoHash"`
	Index     uint32      `gorm:"column:index;not null;<-:create" json:"index"`
	Path      string      `gorm:"column:path;primaryKey;<-:create" json:"path"`
	Extension NullString  `gorm:"column:extension;<-:false" json:"extension"`
	Size      uint64      `gorm:"column:size;not null" json:"size"`
	CreatedAt time.Time   `gorm:"column:created_at;not null;<-:create" json:"createdAt"`
	UpdatedAt time.Time   `gorm:"column:updated_at;not null" json:"updatedAt"`
}

// TableName TorrentFile's table name
func (*TorrentFile) TableName() string {
	return TableNameTorrentFile
}
