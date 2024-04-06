package tmdb

import (
	"context"
	"github.com/go-resty/resty/v2"
)

type Requester interface {
	Request(ctx context.Context, path string, queryParams map[string]string, result interface{}) (*resty.Response, error)
}

type requester struct {
	resty *resty.Client
}

func (r requester) Request(ctx context.Context, path string, queryParams map[string]string, result interface{}) (*resty.Response, error) {
	res, err := r.resty.R().
		SetContext(ctx).
		SetQueryParams(queryParams).
		SetResult(&result).
		Get(path)
	if err == nil {
		if !res.IsSuccess() {
			if res.StatusCode() == 401 {
				err = ErrUnauthorized
			} else if res.StatusCode() == 404 {
				err = ErrNotFound
			} else {
				err = newError(res.Status())
			}
		}
	}
	return res, err
}
