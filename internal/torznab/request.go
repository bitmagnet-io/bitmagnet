package torznab

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type SearchRequest struct {
	Query          string
	Type           string
	Cats           []int
	ImdbId         model.NullString
	TmdbId         model.NullString
	Season         model.NullInt
	Episode        model.NullInt
	Attrs          []string
	Extended       bool
	Limit          model.NullUint
	Offset         model.NullUint
	OrderBy        search.TorrentContentOrderBy
	OrderDirection search.OrderDirection
	Tags           []string
	PermaLinkBase  string
}
