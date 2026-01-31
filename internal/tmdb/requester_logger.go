package tmdb

import (
	"context"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type requesterLogger struct {
	requester Requester
	logger    *zap.Logger
}

func (r requesterLogger) Request(
	ctx context.Context,
	path string,
	queryParams map[string]string,
	result any,
) (*resty.Response, error) {
	res, err := r.requester.Request(ctx, path, queryParams, result)

	fields := []zap.Field{
		zap.String("path", path),
		zap.Any("params", queryParams),
	}

	if res != nil {
		fields = append(fields, zap.String("status", res.Status()), zap.Any("trace", res.Request.TraceInfo()))
	}

	if err == nil {
		r.logger.Debug("request succeeded", fields...)
	} else {
		fields = append(fields, zap.Error(err))
		r.logger.Error("request failed", fields...)
	}

	return res, err
}
