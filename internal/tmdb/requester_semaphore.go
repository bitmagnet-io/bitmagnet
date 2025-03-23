package tmdb

import (
	"context"
	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/semaphore"
)

type requesterSemaphore struct {
	requester Requester
	semaphore *semaphore.Weighted
}

func (r requesterSemaphore) Request(
	ctx context.Context,
	path string,
	queryParams map[string]string,
	result interface{},
) (*resty.Response, error) {
	if err := r.semaphore.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer r.semaphore.Release(1)
	return r.requester.Request(ctx, path, queryParams, result)
}
