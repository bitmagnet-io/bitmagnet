package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (f *TorrentHint) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{
		UpdateAll: true,
	})
	return nil
}

func (h TorrentHint) IsNil() bool {
	return h.ContentType.IsNil()
}

func (h TorrentHint) NullContentType() NullContentType {
	if h.IsNil() {
		return NullContentType{}
	}
	return NewNullContentType(h.ContentType)
}

func (h TorrentHint) ContentRef() Maybe[ContentRef] {
	if h.ContentID.Valid {
		return MaybeValid(ContentRef{
			Type:   h.ContentType,
			Source: h.ContentSource.String,
			ID:     h.ContentID.String,
		})
	}
	return Maybe[ContentRef]{}
}
