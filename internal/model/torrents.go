package model

import (
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/lexer"
	"github.com/facette/natsort"
	"gorm.io/gorm"
)

func (t *Torrent) AfterFind(_ *gorm.DB) error {
	if t.Files != nil {
		sort.Slice(t.Files, func(i, j int) bool {
			return t.Files[i].Path < t.Files[j].Path
		})
	}

	if t.Tags != nil {
		sort.Slice(t.Tags, func(i, j int) bool {
			return natsort.Compare(t.Tags[i].Name, t.Tags[j].Name)
		})
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

var cutoff = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

func (t Torrent) PublishedAt() time.Time {
	publishedAt := t.CreatedAt

	for _, source := range t.Sources {
		dt := source.CreatedAt
		if source.PublishedAt.Valid && source.PublishedAt.Time.After(cutoff) {
			dt = source.PublishedAt.Time
		}

		if dt.Before(publishedAt) {
			publishedAt = dt
		}
	}

	return publishedAt
}

func (t Torrent) MagnetURI() string {
	return "magnet:?xt=urn:btih:" + t.InfoHash.String() +
		"&dn=" + url.QueryEscape(t.Name) +
		"&xl=" + strconv.FormatUint(uint64(t.Size), 10)
}

// HasFilesInfo returns true if we know about the files in this torrent.
func (t Torrent) HasFilesInfo() bool {
	return t.FilesStatus == FilesStatusSingle || t.FilesStatus == FilesStatusMulti || len(t.Files) > 0
}

func (t Torrent) SingleFile() bool {
	return t.FilesStatus == FilesStatusSingle
}

func (t Torrent) BaseName() string {
	baseName := t.Name
	if t.Extension.Valid {
		baseName = baseName[:len(baseName)-len(t.Extension.String)-1]
	}

	return baseName
}

func (t Torrent) FileExtensions() []string {
	switch t.FilesStatus {
	case FilesStatusSingle:
		exts := make([]string, 0, 1)
		ext := FileExtensionFromPath(t.Name)

		if ext.Valid {
			exts = append(exts, ext.String)
		}

		return exts
	default:
		exts := make([]string, 0, len(t.Files))
		extMap := make(map[string]struct{})

		for _, file := range t.Files {
			ext := FileExtensionFromPath(file.Path)
			if ext.Valid {
				if _, ok := extMap[ext.String]; !ok {
					extMap[ext.String] = struct{}{}

					exts = append(exts, ext.String)
				}
			}
		}

		return exts
	}
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
	for _, thisFt := range t.FileTypes() {
		for _, ft := range fts {
			if ft == thisFt {
				return NewNullBool(true)
			}
		}
	}

	return NewNullBool(false)
}

func (t Torrent) TagNames() []string {
	tagNames := make([]string, 0, len(t.Tags))
	for _, tag := range t.Tags {
		tagNames = append(tagNames, tag.Name)
	}

	return tagNames
}

// fileSearchStrings returns a list of strings extracted from file paths, for inclusion in the text search vector.
// To reduce duplication, common prefixes and suffixes are deduplicated.
func (t Torrent) fileSearchStrings() []string {
	firstPass := make([]string, 0, len(t.Files))

	var prevPath string
outer:
	for _, f := range t.Files {
		i := 0
		for {
			if i >= len(f.Path) {
				continue outer
			}
			if i >= len(prevPath) || prevPath[i] != f.Path[i] {
				break
			}
			i++
		}
		for i != 0 && lexer.IsWordChar(rune(f.Path[i])) {
			i--
		}
		firstPass = append(firstPass, f.Path[i:])
		prevPath = f.Path
	}

	searchStrings := make([]string, 0, len(firstPass))

	for i := range firstPass {
		longestSuffixLength := 0

		for j := range i {
			l := 0

			for l < len(firstPass[i]) &&
				l < len(firstPass[j]) &&
				firstPass[i][len(firstPass[i])-l-1] == firstPass[j][len(firstPass[j])-l-1] {
				l++
			}

			if l > longestSuffixLength {
				longestSuffixLength = l
			}
		}

		for longestSuffixLength != 0 &&
			lexer.IsWordChar(rune(firstPass[i][len(firstPass[i])-longestSuffixLength])) {
			longestSuffixLength--
		}

		str := strings.TrimSpace(firstPass[i][:len(firstPass[i])-longestSuffixLength])
		if str != "" {
			searchStrings = append(searchStrings, str)
		}
	}

	return searchStrings
}
