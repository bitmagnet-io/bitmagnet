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

	bestMatch, ok := levenshteinFindBestMatch[tmdb.SearchMovieResult](
		title,
		searchResult.Results,
		func(item tmdb.SearchMovieResult) []string {
			return []string{item.Title, item.OriginalTitle}
		},
	)

	if !ok {
		return model.Content{}, classification.ErrUnmatched
	}

	return c.tmdbGetMovieByTMDBID(bestMatch.ID)
}

func (c executionContext) tmdbSearchTVShow(title string, year model.Year) (model.Content, error) {
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

	bestMatch, ok := levenshteinFindBestMatch[tmdb.SearchTvResult](
		title,
		searchResult.Results,
		func(item tmdb.SearchTvResult) []string {
			return []string{item.Name, item.OriginalName}
		},
	)

	if !ok {
		return model.Content{}, classification.ErrUnmatched
	}

	return c.tmdbGetTVShowByTMDBID(bestMatch.ID)
}

func (c executionContext) tmdbGetMovieByTMDBID(id int64) (movie model.Content, err error) {
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

func (c executionContext) tmdbGetTVShowByTMDBID(id int64) (movie model.Content, err error) {
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

func (c executionContext) tmdbGetTMDBIDByExternalID(ref model.ContentRef) (int64, error) {
	externalSource, externalID, externalSourceErr := tmdb.ExternalSource(ref)
	if externalSourceErr != nil {
		return 0, externalSourceErr
	}
	byIDResult, byIDErr := c.tmdbClient.FindByID(c.Context, tmdb.FindByIDRequest{
		ExternalSource: externalSource,
		ExternalID:     externalID,
	})
	if byIDErr != nil {
		return 0, byIDErr
	}
	switch ref.Type {
	case model.ContentTypeMovie, model.ContentTypeXxx:
		if len(byIDResult.MovieResults) == 0 {
			return 0, classification.ErrUnmatched
		}
		return byIDResult.MovieResults[0].ID, nil
	case model.ContentTypeTvShow:
		if len(byIDResult.TvResults) == 0 {
			return 0, classification.ErrUnmatched
		}
		return byIDResult.TvResults[0].ID, nil
	default:
		return 0, classification.ErrUnmatched
	}
}
