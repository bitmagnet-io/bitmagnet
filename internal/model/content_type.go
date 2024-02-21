package model

// ContentType represents the type of content
// ENUM(movie, tv_show, music, ebook, comic, audiobook, game, software, xxx)
type ContentType string

func (c ContentType) Label() string {
	return c.String()
}

func (c ContentType) IsNil() bool {
	return c == ""
}

func (c ContentType) IsVideo() bool {
	return c == ContentTypeMovie || c == ContentTypeTvShow || c == ContentTypeXxx
}

// A map of file extensions to associated content types.
// The map includes only extensions where the content type can be very reliably inferred.
var extensionToContentTypeMap = map[string]ContentType{
	"m4b":     ContentTypeAudiobook,
	"cb7":     ContentTypeComic,
	"cba":     ContentTypeComic,
	"cbr":     ContentTypeComic,
	"cbt":     ContentTypeComic,
	"cbz":     ContentTypeComic,
	"epub":    ContentTypeEbook,
	"mobi":    ContentTypeEbook,
	"azw":     ContentTypeEbook,
	"azw3":    ContentTypeEbook,
	"pdf":     ContentTypeEbook,
	"chm":     ContentTypeEbook,
	"doc":     ContentTypeEbook,
	"docx":    ContentTypeEbook,
	"odt":     ContentTypeEbook,
	"rtf":     ContentTypeEbook,
	"djvu":    ContentTypeEbook,
	"ape":     ContentTypeMusic,
	"flac":    ContentTypeMusic,
	"exe":     ContentTypeSoftware,
	"dmg":     ContentTypeSoftware,
	"app":     ContentTypeSoftware,
	"apk":     ContentTypeSoftware,
	"deb":     ContentTypeSoftware,
	"rpm":     ContentTypeSoftware,
	"jar":     ContentTypeSoftware,
	"dll":     ContentTypeSoftware,
	"lua":     ContentTypeSoftware,
	"package": ContentTypeSoftware,
	"pkg":     ContentTypeSoftware,
}

func ContentTypeFromExtension(ext string) NullContentType {
	ct, ok := extensionToContentTypeMap[ext]
	if !ok {
		return NullContentType{}
	}
	return NewNullContentType(ct)
}
