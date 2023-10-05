package model

import (
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
	"gorm.io/gorm"
)

func (t *Torrent) BeforeCreate(tx *gorm.DB) error {
	if len(t.Contents) == 0 {
		t.Contents = []TorrentContent{
			{
				InfoHash:     t.InfoHash,
				Title:        t.Name,
				SearchString: regex.NormalizeString(t.Name),
			},
		}
	}
	return nil
}

// Seeders returns the highest number of seeders from all sources
// todo: Add up bloom filters
func (t Torrent) Seeders() NullUint {
	seeders := NullUint{}
	for _, source := range t.Sources {
		if source.Seeders.Valid {
			seeders.Valid = true
			if source.Seeders.Uint > seeders.Uint {
				seeders.Uint = source.Seeders.Uint
			}
		}
	}
	return seeders
}

// Leechers returns the highest number of leechers from all sources
func (t Torrent) Leechers() NullUint {
	leechers := NullUint{}
	for _, source := range t.Sources {
		if source.Leechers.Valid {
			leechers.Valid = true
			if source.Leechers.Uint > leechers.Uint {
				leechers.Uint = source.Leechers.Uint
			}
		}
	}
	return leechers
}

func (t Torrent) Magnet() string {
	return "magnet:?xt=urn:btih:" + t.InfoHash.String()
}

// HasFilesInfo returns true if we know about the files in this torrent.
func (t Torrent) HasFilesInfo() bool {
	return t.FilesStatus == FilesStatusSingle || t.FilesStatus == FilesStatusMulti
}

func (t Torrent) WantFilesInfo() bool {
	return t.FilesStatus == FilesStatusNoInfo
}

func (t Torrent) SingleFile() bool {
	return t.FilesStatus == FilesStatusSingle
}

func (t Torrent) FileExtensions() []string {
	switch t.FilesStatus {
	case FilesStatusSingle:
		exts := make([]string, 0, 1)
		ext := fileExtensionFromPath(t.Name)
		if ext.Valid {
			exts = append(exts, ext.String)
		}
		return exts
	case FilesStatusMulti:
		exts := make([]string, 0, len(t.Files))
		extMap := make(map[string]struct{})
		for _, file := range t.Files {
			ext := fileExtensionFromPath(file.Path)
			if ext.Valid {
				if _, ok := extMap[ext.String]; !ok {
					extMap[ext.String] = struct{}{}
					exts = append(exts, ext.String)
				}
			}
		}
		return exts
	}
	return nil
}

func (t Torrent) FileType() NullFileType {
	if t.Extension.Valid {
		return FileTypeFromExtension(t.Extension.String)
	}
	return NullFileType{}
}

func (t Torrent) FileTypes() []FileType {
	exts := t.FileExtensions()
	typesMap := make(map[FileType]struct{})
	types := make([]FileType, 0, len(exts))
	for _, ext := range exts {
		if ft := FileTypeFromExtension(ext); ft.Valid {
			if _, ok := typesMap[ft.FileType]; !ok {
				typesMap[ft.FileType] = struct{}{}
				types = append(types, ft.FileType)
			}
		}
	}
	return types
}

func (t Torrent) HasFileType(fts ...FileType) NullBool {
	if !t.HasFilesInfo() {
		return NullBool{}
	}
	for _, thisFt := range t.FileTypes() {
		for _, ft := range fts {
			if ft == thisFt {
				return NewNullBool(true)
			}
		}
	}
	return NewNullBool(false)
}
