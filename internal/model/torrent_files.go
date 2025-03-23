package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"regexp"
	"strings"
)

func (*TorrentFile) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{
		DoNothing: true,
	})
	return nil
}

func (f TorrentFile) BasePath() string {
	baseName := f.Path
	if f.Extension.Valid {
		baseName = baseName[:len(baseName)-len(f.Extension.String)-1]
	}
	return baseName
}

func (f TorrentFile) BaseName() string {
	basePathParts := strings.Split(f.BasePath(), "/")
	return basePathParts[len(basePathParts)-1]
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
