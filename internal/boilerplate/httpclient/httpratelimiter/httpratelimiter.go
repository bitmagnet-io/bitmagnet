package httpratelimiter

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/httpclient"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

func NewDecorator(limit time.Duration, burst int) httpclient.TransportDecorator {
	return func(t http.RoundTripper) http.RoundTripper {
		return &httpRateLimiter{
			limiter:   rate.NewLimiter(rate.Every(limit), burst),
			transport: t,
		}
	}
}

type httpRateLimiter struct {
	limiter   *rate.Limiter
	transport http.RoundTripper
}

func (l *httpRateLimiter) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := l.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}
	return l.transport.RoundTrip(req)
}
