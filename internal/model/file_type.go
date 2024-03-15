package model

import (
	"sort"
	"strings"
)

// FileType represents the general type of a file
/* ENUM(
  archive,
	audio,
  data,
  document,
  image,
  software,
  subtitles,
	video,
) */
type FileType string

var extensionToFileTypeMap = map[string]FileType{
	// Archive
	"zip": FileTypeArchive,
	"rar": FileTypeArchive,
	"tar": FileTypeArchive,
	"gz":  FileTypeArchive,
	"7z":  FileTypeArchive,
	"iso": FileTypeArchive,
	"bz2": FileTypeArchive,

	// Audio
	"mp3":  FileTypeAudio,
	"wav":  FileTypeAudio,
	"flac": FileTypeAudio,
	"aac":  FileTypeAudio,
	"ogg":  FileTypeAudio,
	"m4a":  FileTypeAudio,
	"m4b":  FileTypeAudio,
	"mid":  FileTypeAudio,
	"dsf":  FileTypeAudio,

	// Data
	"csv":  FileTypeData,
	"json": FileTypeData,
	"xml":  FileTypeData,
	"xls":  FileTypeData,
	"xlsx": FileTypeData,

	// Document
	"pdf":  FileTypeDocument,
	"doc":  FileTypeDocument,
	"docx": FileTypeDocument,
	"otf":  FileTypeDocument,
	"ppt":  FileTypeDocument,
	"pptx": FileTypeDocument,
	"html": FileTypeDocument,
	"htm":  FileTypeDocument,
	"epub": FileTypeDocument,
	"mobi": FileTypeDocument,
	"azw":  FileTypeDocument,
	"azw3": FileTypeDocument,
	"rtf":  FileTypeDocument,
	"txt":  FileTypeDocument,
	"md":   FileTypeDocument,
	"nfo":  FileTypeDocument,
	"djvu": FileTypeDocument,

	// Image
	"jpg":  FileTypeImage,
	"jpeg": FileTypeImage,
	"png":  FileTypeImage,
	"gif":  FileTypeImage,
	"bmp":  FileTypeImage,
	"svg":  FileTypeImage,
	"dds":  FileTypeImage,
	"psd":  FileTypeImage,
	"tif":  FileTypeImage,
	"tiff": FileTypeImage,
	"ico":  FileTypeImage,

	// Software
	"exe":     FileTypeSoftware,
	"bin":     FileTypeSoftware,
	"sh":      FileTypeSoftware,
	"bat":     FileTypeSoftware,
	"msi":     FileTypeSoftware,
	"apk":     FileTypeSoftware,
	"dmg":     FileTypeSoftware,
	"pkg":     FileTypeSoftware,
	"deb":     FileTypeSoftware,
	"rpm":     FileTypeSoftware,
	"jar":     FileTypeSoftware,
	"dll":     FileTypeSoftware,
	"lua":     FileTypeSoftware,
	"package": FileTypeSoftware,

	// Subtitles
	"srt": FileTypeSubtitles,
	"sub": FileTypeSubtitles,
	"vtt": FileTypeSubtitles,

	// Video
	"mp4":  FileTypeVideo,
	"mkv":  FileTypeVideo,
	"avi":  FileTypeVideo,
	"mov":  FileTypeVideo,
	"wmv":  FileTypeVideo,
	"flv":  FileTypeVideo,
	"m4v":  FileTypeVideo,
	"mpg":  FileTypeVideo,
	"mpeg": FileTypeVideo,
}

var fileTypeToExtensionsMap map[FileType][]string

func init() {
	m := make(map[FileType][]string)
	for ext, ft := range extensionToFileTypeMap {
		if _, ok := m[ft]; !ok {
			m[ft] = make([]string, 0)
		}
		m[ft] = append(m[ft], ext)
	}
	for ft := range m {
		sort.Strings(m[ft])
	}
	fileTypeToExtensionsMap = m
}

func (ft FileType) Extensions() []string {
	return fileTypeToExtensionsMap[ft]
}

func (ft FileType) Label() string {
	return strings.ToUpper(string(ft[0])) + string(ft[1:])
}

func FileTypeFromExtension(ext string) NullFileType {
	if t, ok := extensionToFileTypeMap[ext]; ok {
		return NullFileType{FileType: t, Valid: true}
	}
	return NullFileType{}
}
