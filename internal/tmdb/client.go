package tmdb

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"strconv"
	"strings"
)

type client struct {
	resty          *resty.Client
	limiter        *rate.Limiter
	isUnauthorized *concurrency.AtomicValue[bool]
	logger         *zap.SugaredLogger
}

func newError(msg string) error {
	return fmt.Errorf("TMDB request failed: %s", msg)
}

var ErrUnauthorized = newError("401 Unauthorized")
var ErrNotFound = newError("404 Not Found")

func (c client) ValidateApiKey(ctx context.Context) error {
	return c.request(ctx, "/authentication", nil, nil)
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
	err := c.request(ctx, "/search/movie", queryParams, &response)
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
	err := c.request(ctx, "/movie/"+strconv.FormatInt(request.ID, 10), queryParams, &response)
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
	err := c.request(ctx, "/search/tv", queryParams, &response)
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
	err := c.request(ctx, "/tv/"+strconv.FormatInt(request.SeriesID, 10), queryParams, &response)
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
	err := c.request(ctx, "/find/"+request.ExternalID, queryParams, &response)
	return response, err
}

func (c client) request(ctx context.Context, path string, queryParams map[string]string, result interface{}) error {
	var (
		res *resty.Response
		err error
	)
	if c.isUnauthorized.Get() {
		err = ErrUnauthorized
	} else {
		err = c.limiter.Wait(ctx)
	}
	if err == nil {
		res, err = c.resty.R().
			SetContext(ctx).
			SetQueryParams(queryParams).
			SetResult(&result).
			Get(path)
	}
	if err == nil {
		if !res.IsSuccess() {
			if res.StatusCode() == 401 {
				c.isUnauthorized.Set(true)
				err = ErrUnauthorized
			} else if res.StatusCode() == 404 {
				err = ErrNotFound
			} else {
				err = newError(res.Status())
			}
		}
	}
	kvs := []interface{}{"path", path, "queryParams", queryParams}
	if res != nil {
		kvs = append(kvs, "status", res.Status(), "trace", res.Request.TraceInfo())
	}
	if err == nil {
		c.logger.Debugw("request succeeded", kvs...)
	} else {
		kvs = append(kvs, "error", err)
		c.logger.Errorw("request failed", kvs...)
	}
	return err
}
