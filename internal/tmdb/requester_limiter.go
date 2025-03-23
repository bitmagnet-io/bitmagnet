package tmdb

import (
	"context"

	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
)

type requesterLimiter struct {
	requester Requester
	limiter   *rate.Limiter
}

func (r requesterLimiter) Request(
	ctx context.Context,
	path string,
	queryParams map[string]string,
	result any,
) (*resty.Response, error) {
	if err := r.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	return r.requester.Request(ctx, path, queryParams, result)
}
