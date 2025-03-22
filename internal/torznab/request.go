package torznab

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type SearchRequest struct {
	Query    string
	Type     string
	Cats     []int
	IMDBID   model.NullString
	TMDBID   model.NullString
	Season   model.NullInt
	Episode  model.NullInt
	Attrs    []string
	Extended bool
	Limit    model.NullUint
	Offset   model.NullUint
}
