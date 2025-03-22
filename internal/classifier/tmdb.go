package classifier

import (
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
)

func (c executionContext) tmdbSearchMovie(title string, year model.Year) (model.Content, error) {
	req := tmdb.SearchMovieRequest{
		Query:        title,
		IncludeAdult: true,
	}
	if !year.IsNil() {
		req.Year = year
	}
	searchResult, searchErr := c.tmdbClient.SearchMovie(c.Context, req)
	if searchErr != nil {
		return model.Content{}, searchErr
	}

	if bestMatch, ok := levenshteinFindBestMatch[tmdb.SearchMovieResult](
		title,
		searchResult.Results,
		func(item tmdb.SearchMovieResult) []string {
			return []string{item.Title, item.OriginalTitle}
		},
	); !ok {
		return model.Content{}, classification.ErrUnmatched
	} else {
		return c.tmdbGetMovieByTmbdId(bestMatch.ID)
	}
}

func (c executionContext) tmdbSearchTvShow(title string, year model.Year) (model.Content, error) {
	req := tmdb.SearchTvRequest{
		Query:        title,
		IncludeAdult: true,
	}
	if !year.IsNil() {
		req.FirstAirDateYear = year
	}
	searchResult, searchErr := c.tmdbClient.SearchTv(c.Context, req)
	if searchErr != nil {
		return model.Content{}, searchErr
	}

	if bestMatch, ok := levenshteinFindBestMatch[tmdb.SearchTvResult](
		title,
		searchResult.Results,
		func(item tmdb.SearchTvResult) []string {
			return []string{item.Name, item.OriginalName}
		},
	); !ok {
		return model.Content{}, classification.ErrUnmatched
	} else {
		return c.tmdbGetTvShowByTmbdId(bestMatch.ID)
	}
}

func (c executionContext) tmdbGetMovieByTmbdId(id int64) (movie model.Content, err error) {
	d, getDetailsErr := c.tmdbClient.MovieDetails(c.Context, tmdb.MovieDetailsRequest{
		ID: id,
	})
	if getDetailsErr != nil {
		if errors.Is(getDetailsErr, tmdb.ErrNotFound) {
			getDetailsErr = classification.ErrUnmatched
		}
		err = getDetailsErr
		return
	}
	return tmdb.MovieDetailsToMovieModel(d)
}

func (c executionContext) tmdbGetTvShowByTmbdId(id int64) (movie model.Content, err error) {
	d, getDetailsErr := c.tmdbClient.TvDetails(c.Context, tmdb.TvDetailsRequest{
		SeriesID:         id,
		AppendToResponse: []string{"external_ids"},
	})
	if getDetailsErr != nil {
		if errors.Is(getDetailsErr, tmdb.ErrNotFound) {
			getDetailsErr = classification.ErrUnmatched
		}
		err = getDetailsErr
		return
	}
	return tmdb.TvShowDetailsToTvShowModel(d)
}

func (c executionContext) tmdbGetTmdbIdByExternalId(ref model.ContentRef) (int64, error) {
	externalSource, externalID, externalSourceErr := tmdb.ExternalSource(ref)
	if externalSourceErr != nil {
		return 0, externalSourceErr
	}
	byIdResult, byIdErr := c.tmdbClient.FindByID(c.Context, tmdb.FindByIDRequest{
		ExternalSource: externalSource,
		ExternalID:     externalID,
	})
	if byIdErr != nil {
		return 0, byIdErr
	}
	switch ref.Type {
	case model.ContentTypeMovie, model.ContentTypeXxx:
		if len(byIdResult.MovieResults) == 0 {
			return 0, classification.ErrUnmatched
		}
		return byIdResult.MovieResults[0].ID, nil
	case model.ContentTypeTvShow:
		if len(byIdResult.TvResults) == 0 {
			return 0, classification.ErrUnmatched
		}
		return byIdResult.TvResults[0].ID, nil
	default:
		return 0, classification.ErrUnmatched
	}
}
