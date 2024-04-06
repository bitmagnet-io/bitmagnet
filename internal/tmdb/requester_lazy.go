package tmdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type requesterLazy struct {
	once      sync.Once
	config    Config
	logger    *zap.SugaredLogger
	err       error
	requester Requester
}

func (r *requesterLazy) Request(ctx context.Context, path string, queryParams map[string]string, result interface{}) (*resty.Response, error) {
	r.once.Do(func() {
		r.requester, r.err = newRequester(ctx, r.config, r.logger)
	})
	if r.err != nil {
		return nil, r.err
	}
	return r.requester.Request(ctx, path, queryParams, result)
}

func newRequester(ctx context.Context, config Config, logger *zap.SugaredLogger) (Requester, error) {
	if config.ApiKey == defaultTmdbApiKey {
		logger.Warnln("you are using the default TMDB api key; TMDB requests will be limited to 1 per second; to remove this warning please configure a personal TMDB api key")
		config.RateLimit = time.Second
		config.RateLimitBurst = 1
	}
	r := requesterLogger{
		requester: requesterFailFast{
			requester: requesterLimiter{
				requester: requester{
					resty: resty.New().
						SetBaseURL(config.BaseUrl).
						SetQueryParam("api_key", config.ApiKey).
						SetRetryCount(3).
						SetRetryWaitTime(2 * time.Second).
						SetRetryMaxWaitTime(20 * time.Second).
						SetTimeout(10 * time.Second).
						EnableTrace().
						SetLogger(logger),
				},
				limiter: rate.NewLimiter(rate.Every(config.RateLimit), config.RateLimitBurst),
			},
			isUnauthorized: &concurrency.AtomicValue[bool]{},
		},
		logger: logger,
	}
	err := client{r}.ValidateApiKey(context.Background())
	if errors.Is(err, ErrUnauthorized) {
		if config.ApiKey == defaultTmdbApiKey {
			return r, fmt.Errorf("default api key is invalid: %w", err)
		}
		logger.Errorw("invalid api key, falling back to default", "error", err)
		config.ApiKey = defaultTmdbApiKey
		return newRequester(ctx, config, logger)
	}
	return r, err
}
