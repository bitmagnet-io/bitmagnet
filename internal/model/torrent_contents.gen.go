// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database/fts"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

const TableNameTorrentContent = "torrent_contents"

// TorrentContent mapped from table <torrent_contents>
type TorrentContent struct {
	ID              string              `gorm:"column:id;primaryKey;<-:false" json:"id"`
	InfoHash        protocol.ID         `gorm:"column:info_hash;not null;<-:create" json:"infoHash"`
	ContentType     NullContentType     `gorm:"column:content_type" json:"contentType"`
	ContentSource   NullString          `gorm:"column:content_source" json:"contentSource"`
	ContentID       NullString          `gorm:"column:content_id" json:"contentId"`
	Languages       Languages           `gorm:"column:languages;serializer:json" json:"languages"`
	Episodes        Episodes            `gorm:"column:episodes;serializer:json" json:"episodes"`
	VideoResolution NullVideoResolution `gorm:"column:video_resolution" json:"videoResolution"`
	VideoSource     NullVideoSource     `gorm:"column:video_source" json:"videoSource"`
	VideoCodec      NullVideoCodec      `gorm:"column:video_codec" json:"videoCodec"`
	Video3D         NullVideo3D         `gorm:"column:video_3d" json:"video3D"`
	VideoModifier   NullVideoModifier   `gorm:"column:video_modifier" json:"videoModifier"`
	ReleaseGroup    NullString          `gorm:"column:release_group" json:"releaseGroup"`
	CreatedAt       time.Time           `gorm:"column:created_at;not null;<-:create" json:"createdAt"`
	UpdatedAt       time.Time           `gorm:"column:updated_at;not null" json:"updatedAt"`
	Tsv             fts.Tsvector        `gorm:"column:tsv" json:"tsv"`
	Seeders         NullUint            `gorm:"column:seeders" json:"seeders"`
	Leechers        NullUint            `gorm:"column:leechers" json:"leechers"`
	PublishedAt     time.Time           `gorm:"column:published_at;not null;default:1999-01-01 00:00:00+00" json:"publishedAt"`
	Size            uint                `gorm:"column:size;not null" json:"size"`
	FilesCount      NullUint            `gorm:"column:files_count" json:"filesCount"`
	Torrent         Torrent             `gorm:"foreignKey:InfoHash;references:InfoHash" json:"torrent"`
	Content         Content             `gorm:"foreignKey:ContentType,ContentSource,ContentID;references:Type,Source,ID" json:"content"`
}

// TableName TorrentContent's table name
func (*TorrentContent) TableName() string {
	return TableNameTorrentContent
}
