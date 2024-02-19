package model

// ContentType represents the type of content
// ENUM(movie, tv_show, music, ebook, audiobook, game, software, xxx)
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

var extensionToContentTypeMap = map[string]ContentType{
	"m4b":     ContentTypeAudiobook,
	"epub":    ContentTypeEbook,
	"mobi":    ContentTypeEbook,
	"azw":     ContentTypeEbook,
	"azw3":    ContentTypeEbook,
	"pdf":     ContentTypeEbook,
	"cbr":     ContentTypeEbook,
	"cbz":     ContentTypeEbook,
	"cb7":     ContentTypeEbook,
	"cbt":     ContentTypeEbook,
	"cba":     ContentTypeEbook,
	"chm":     ContentTypeEbook,
	"doc":     ContentTypeEbook,
	"docx":    ContentTypeEbook,
	"odt":     ContentTypeEbook,
	"rtf":     ContentTypeEbook,
	"djvu":    ContentTypeEbook,
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
