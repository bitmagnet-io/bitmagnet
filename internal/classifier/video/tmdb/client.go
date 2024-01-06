package tmdb

import (
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/cyruzin/golang-tmdb"
)

type Client interface {
	MovieClient
	TvShowClient
}

type client struct {
	c *tmdb.Client
	s search.Search
}

const SourceTmdb = "tmdb"

var (
	ErrNotFound      = errors.New("not found")
	ErrUnknownSource = errors.New("unknown source")
)
