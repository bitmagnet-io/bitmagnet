package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const attachTmdbContentBySearchName = "attach_tmdb_content_by_search"

type attachTmdbContentBySearchAction struct{}

func (attachTmdbContentBySearchAction) name() string {
	return attachTmdbContentBySearchName
}

var attachTmdbContentBySearchPayloadSpec = payloadLiteral[string]{
	literal:     attachTmdbContentBySearchName,
	description: "Attempt to attach content from the TMDB API with a search on the torrent name",
}

func (a attachTmdbContentBySearchAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := attachTmdbContentBySearchPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			if !cl.BaseTitle.Valid {
				return cl, classification.ErrUnmatched
			}
			var content *model.Content
			switch cl.ContentType.ContentType {
			case model.ContentTypeTvShow:
				if result, searchErr := ctx.tmdb_searchTvShow(cl.BaseTitle.String, cl.Date.Year); searchErr != nil {
					return cl, searchErr
				} else {
					content = &result
				}
			default:
				if len(cl.Episodes) > 0 {
					return cl, classification.ErrUnmatched
				}
				if result, searchErr := ctx.tmdb_searchMovie(cl.BaseTitle.String, cl.Date.Year); searchErr != nil {
					return cl, searchErr
				} else {
					content = &result
				}
			}
			cl.AttachContent(content)
			return cl, nil
		},
	}, nil
}

func (attachTmdbContentBySearchAction) JsonSchema() JsonSchema {
	return attachTmdbContentBySearchPayloadSpec.JsonSchema()
}
