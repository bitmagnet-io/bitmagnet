package httplogger

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpclient"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func NewDecorator(logger *zap.SugaredLogger) httpclient.TransportDecorator {
	return func(t http.RoundTripper) http.RoundTripper {
		return &httpLogger{
			logger:    logger,
			transport: t,
		}
	}
}

type httpLogger struct {
	logger    *zap.SugaredLogger
	transport http.RoundTripper
}

func (l *httpLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	l.logger.Debugw("http request", "method", req.Method, "url", req.URL.String())
	startTime := time.Now()
	resp, err := l.transport.RoundTrip(req)
	duration := time.Since(startTime)
	if err != nil {
		l.logger.Errorw("http request failed", "method", req.Method, "url", req.URL.String(), "error", err, "duration", duration)
		return nil, err
	}
	l.logger.Debugw("http response", "method", req.Method, "url", req.URL.String(), "status_code", resp.StatusCode, "duration", duration)
	return resp, nil
}
