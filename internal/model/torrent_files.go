package model

import (
	"regexp"
	"strings"
)

var fileExtensionRegex = regexp.MustCompile(`[^/.]\.([a-z0-9]+)$`)

func fileExtensionFromPath(path string) NullString {
	match := fileExtensionRegex.FindStringSubmatch(strings.ToLower(path))
	if len(match) == 2 {
		return NewNullString(match[1])
	}
	return NullString{}
}

func fileTypeFromPath(path string) NullFileType {
	extension := fileExtensionFromPath(path)
	if extension.Valid {
		return FileTypeFromExtension(extension.String)
	}
	return NullFileType{}
}

func (f TorrentFile) FileType() NullFileType {
	return fileTypeFromPath(f.Path)
}
