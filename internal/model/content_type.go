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
