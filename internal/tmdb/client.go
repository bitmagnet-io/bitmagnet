package tmdb

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type client struct {
	requester Requester
}

func newError(msg string) error {
	return fmt.Errorf("TMDB request failed: %s", msg)
}

var ErrUnauthorized = newError("401 Unauthorized")
var ErrNotFound = newError("404 Not Found")

func (c client) ValidateApiKey(ctx context.Context) error {
	_, err := c.requester.Request(ctx, "/authentication", nil, nil)
	return err
}

func (c client) SearchMovie(ctx context.Context, request SearchMovieRequest) (SearchMovieResponse, error) {
	queryParams := map[string]string{
		"query": request.Query,
	}
	if request.IncludeAdult {
		queryParams["include_adult"] = "true"
	}
	if request.Language.Valid {
		queryParams["language"] = request.Language.String
	}
	if !request.PrimaryReleaseYear.IsNil() {
		queryParams["primary_release_year"] = request.PrimaryReleaseYear.String()
	}
	if !request.Year.IsNil() {
		queryParams["year"] = request.Year.String()
	}
	if request.Region.Valid {
		queryParams["region"] = request.Region.String
	}
	var response SearchMovieResponse
	_, err := c.requester.Request(ctx, "/search/movie", queryParams, &response)
	return response, err
}

func (c client) MovieDetails(ctx context.Context, request MovieDetailsRequest) (MovieDetailsResponse, error) {
	queryParams := make(map[string]string)
	if len(request.AppendToResponse) > 0 {
		queryParams["append_to_response"] = strings.Join(request.AppendToResponse, ",")
	}
	if request.Language.Valid {
		queryParams["language"] = request.Language.String
	}
	var response MovieDetailsResponse
	_, err := c.requester.Request(ctx, "/movie/"+strconv.FormatInt(request.ID, 10), queryParams, &response)
	return response, err
}

func (c client) SearchTv(ctx context.Context, request SearchTvRequest) (SearchTvResponse, error) {
	queryParams := map[string]string{
		"query": request.Query,
	}
	if !request.FirstAirDateYear.IsNil() {
		queryParams["first_air_date_year"] = request.FirstAirDateYear.String()
	}
	if request.IncludeAdult {
		queryParams["include_adult"] = "true"
	}
	if request.Language.Valid {
		queryParams["language"] = request.Language.String
	}
	var response SearchTvResponse
	_, err := c.requester.Request(ctx, "/search/tv", queryParams, &response)
	return response, err
}

func (c client) TvDetails(ctx context.Context, request TvDetailsRequest) (TvDetailsResponse, error) {
	queryParams := make(map[string]string)
	if len(request.AppendToResponse) > 0 {
		queryParams["append_to_response"] = strings.Join(request.AppendToResponse, ",")
	}
	if request.Language.Valid {
		queryParams["language"] = request.Language.String
	}
	var response TvDetailsResponse
	_, err := c.requester.Request(ctx, "/tv/"+strconv.FormatInt(request.SeriesID, 10), queryParams, &response)
	return response, err
}

func (c client) FindByID(ctx context.Context, request FindByIDRequest) (FindByIDResponse, error) {
	queryParams := map[string]string{
		"external_source": request.ExternalSource,
	}
	if request.Language.Valid {
		queryParams["language"] = request.Language.String
	}
	var response FindByIDResponse
	_, err := c.requester.Request(ctx, "/find/"+request.ExternalID, queryParams, &response)
	return response, err
}
