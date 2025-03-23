package tmdb

import (
	"context"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type requesterLogger struct {
	requester Requester
	logger    *zap.SugaredLogger
}

func (r requesterLogger) Request(
	ctx context.Context,
	path string,
	queryParams map[string]string,
	result any,
) (*resty.Response, error) {
	res, err := r.requester.Request(ctx, path, queryParams, result)
	kvs := []interface{}{"path", path, "queryParams", queryParams}
	if res != nil {
		kvs = append(kvs, "status", res.Status(), "trace", res.Request.TraceInfo())
	}
	if err == nil {
		r.logger.Debugw("request succeeded", kvs...)
	} else {
		kvs = append(kvs, "error", err)
		r.logger.Errorw("request failed", kvs...)
	}
	return res, err
}
