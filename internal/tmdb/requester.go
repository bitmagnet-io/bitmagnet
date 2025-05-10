package tmdb

import (
	"context"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Requester interface {
	Request(ctx context.Context, path string, queryParams map[string]string, result any) (*resty.Response, error)
}

type requester struct {
	resty *resty.Client
}

func (r requester) Request(
	ctx context.Context,
	path string,
	queryParams map[string]string,
	result any,
) (*resty.Response, error) {
	res, err := r.resty.R().
		SetContext(ctx).
		SetQueryParams(queryParams).
		SetResult(&result).
		Get(path)
	if err == nil && !res.IsSuccess() {
		switch res.StatusCode() {
		case http.StatusUnauthorized:
			err = ErrUnauthorized
		case http.StatusNotFound:
			err = ErrNotFound
		default:
			err = newError(res.Status())
		}
	}

	return res, err
}
