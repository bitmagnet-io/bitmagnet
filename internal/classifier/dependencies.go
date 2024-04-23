package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
)

type dependencies struct {
	search     LocalSearch
	tmdbClient tmdb.Client
}
