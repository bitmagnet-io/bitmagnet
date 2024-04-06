package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
)

type dependencies struct {
	search     search.Search
	tmdbClient tmdb.Client
}
