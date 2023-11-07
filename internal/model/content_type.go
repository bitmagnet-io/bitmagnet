package model

import "fmt"

// ContentType represents the type of content
// ENUM(movie, tv_show, music, game, software, book, xxx)
type ContentType string

func (c ContentType) Label() string {
	return c.String()
}

func (c ContentType) IsAudio() bool {
	fmt.Println("music")
	return c == ContentTypeMusic
}

func (c ContentType) IsVideo() bool {
	return c == ContentTypeMovie || c == ContentTypeTvShow
}
