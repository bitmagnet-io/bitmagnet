package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const attachTmdbContentBySearchName = "attach_tmdb_content_by_search"

type attachTmdbContentBySearchAction struct{}

func (attachTmdbContentBySearchAction) name() string {
	return attachTmdbContentBySearchName
}

var attachTmdbContentBySearchSpec = json_spec.Literal[string]{
	Literal:     attachTmdbContentBySearchName,
	Description: "Attempt to attach content from the TMDB API with a search on the torrent name",
}

func (attachTmdbContentBySearchAction) compile(ctx compilerContext) (action, error) {
	if _, err := attachTmdbContentBySearchSpec.Parse(ctx.jsonSpec); err != nil {
		return action{}, ctx.Error(err)
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

func (attachTmdbContentBySearchAction) JSONSchema() json_schema.JSONSchema {
	return attachTmdbContentBySearchSpec.JSONSchema()
}
