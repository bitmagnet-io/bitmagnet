package classifier

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
)

const attachTmdbContentBySearchName = "attach_tmdb_content_by_search"

type attachTmdbContentBySearchAction struct {
	tmdbAction
}

func (attachTmdbContentBySearchAction) Name() string {
	return attachTmdbContentBySearchName
}

var attachTmdbContentBySearchPayloadSpec = payloadLiteral[string]{attachTmdbContentBySearchName}

func (a attachTmdbContentBySearchAction) compileAction(ctx compilerContext) (action, error) {
	if _, err := attachTmdbContentBySearchPayloadSpec.Unmarshal(ctx); err != nil {
		return action{}, ctx.error(err)
	}
	return action{
		run: func(ctx executionContext) (classification.Result, error) {
			cl := ctx.result
			if !cl.BaseTitle.Valid {
				return cl, classification.ErrNoMatch
			}
			var content *model.Content
			switch cl.ContentType.ContentType {
			case model.ContentTypeTvShow:
				if result, searchErr := a.searchTvShow(ctx, cl.BaseTitle.String, cl.Date.Year); searchErr != nil {
					return cl, searchErr
				} else {
					content = &result
				}
			default:
				if len(cl.Episodes) > 0 {
					return cl, classification.ErrNoMatch
				}
				if result, searchErr := a.searchMovie(ctx, cl.BaseTitle.String, cl.Date.Year); searchErr != nil {
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

func (a attachTmdbContentBySearchAction) searchMovie(ctx context.Context, title string, year model.Year) (model.Content, error) {
	req := tmdb.SearchMovieRequest{
		Query:        title,
		IncludeAdult: true,
	}
	if !year.IsNil() {
		req.Year = year
	}
	searchResult, searchErr := a.client.SearchMovie(ctx, req)
	if searchErr != nil {
		return model.Content{}, searchErr
	}
	for _, item := range searchResult.Results {
		if levenshteinCheck(title, []string{item.Title, item.OriginalTitle}, levenshteinThreshold) {
			return a.getMovieByTmbdId(ctx, item.ID)
		}
	}
	return model.Content{}, classification.ErrNoMatch
}

func (a attachTmdbContentBySearchAction) searchTvShow(ctx context.Context, title string, year model.Year) (model.Content, error) {
	req := tmdb.SearchTvRequest{
		Query:        title,
		IncludeAdult: true,
	}
	if !year.IsNil() {
		req.FirstAirDateYear = year
	}
	searchResult, searchErr := a.client.SearchTv(ctx, req)
	if searchErr != nil {
		return model.Content{}, searchErr
	}
	for _, item := range searchResult.Results {
		if levenshteinCheck(title, []string{item.Name, item.OriginalName}, levenshteinThreshold) {
			return a.getTvShowByTmbdId(ctx, item.ID)
		}
	}
	return model.Content{}, classification.ErrNoMatch
}
