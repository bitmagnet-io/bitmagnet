package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"regexp"
	"strings"
)

func (f *TorrentFile) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{
		DoNothing: true,
	})
	return nil
}

var fileExtensionRegex = regexp.MustCompile(`[^/.]\.([a-z0-9]+)$`)

func FileExtensionFromPath(path string) NullString {
	match := fileExtensionRegex.FindStringSubmatch(strings.ToLower(path))
	if len(match) == 2 {
		return NewNullString(match[1])
	}
	return NullString{}
}

func fileTypeFromPath(path string) NullFileType {
	extension := FileExtensionFromPath(path)
	if extension.Valid {
		return FileTypeFromExtension(extension.String)
	}
	return NullFileType{}
}

func (f TorrentFile) FileType() NullFileType {
	return fileTypeFromPath(f.Path)
}
