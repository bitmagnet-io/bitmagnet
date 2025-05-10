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

func (attachTmdbContentBySearchAction) compileAction(ctx compilerContext) (action, error) {
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
				result, searchErr := ctx.tmdbSearchTVShow(cl.BaseTitle.String, cl.Date.Year)
				if searchErr != nil {
					return cl, searchErr
				}
				content = &result
			default:
				if len(cl.Episodes) > 0 {
					return cl, classification.ErrUnmatched
				}
				result, searchErr := ctx.tmdbSearchMovie(cl.BaseTitle.String, cl.Date.Year)
				if searchErr != nil {
					return cl, searchErr
				}
				content = &result
			}
			cl.AttachContent(content)
			return cl, nil
		},
	}, nil
}

func (attachTmdbContentBySearchAction) JSONSchema() JSONSchema {
	return attachTmdbContentBySearchPayloadSpec.JSONSchema()
}
