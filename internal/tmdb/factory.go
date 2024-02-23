package tmdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/go-resty/resty/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"time"
)

type Params struct {
	fx.In
	Config Config
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Client lazy.Lazy[Client]
}

func New(p Params) Result {
	return Result{
		Client: lazy.New(func() (Client, error) {
			return newClient(p.Config, p.Logger.Named("tmdb_client"))
		}),
	}
}

func newClient(config Config, logger *zap.SugaredLogger) (Client, error) {
	if config.ApiKey == defaultTmdbApiKey {
		logger.Warnln("you are using the default TMDB api key; TMDB requests will be limited to 1 per second; to remove this warning please configure a personal TMDB api key")
		config.RateLimit = time.Second
		config.RateLimitBurst = 1
	}
	cl := client{
		resty: resty.New().
			SetBaseURL("https://api.themoviedb.org/3").
			SetQueryParam("api_key", config.ApiKey).
			SetRetryCount(3).
			SetRetryWaitTime(2 * time.Second).
			SetRetryMaxWaitTime(20 * time.Second).
			SetTimeout(10 * time.Second).
			EnableTrace().
			SetLogger(logger),
		limiter:        rate.NewLimiter(rate.Every(config.RateLimit), config.RateLimitBurst),
		isUnauthorized: &concurrency.AtomicValue[bool]{},
		logger:         logger,
	}
	err := cl.ValidateApiKey(context.Background())
	if errors.Is(err, ErrUnauthorized) {
		if config.ApiKey == defaultTmdbApiKey {
			return cl, fmt.Errorf("default api key is invalid: %w", err)
		}
		logger.Errorw("invalid api key, falling back to default", "error", err)
		config.ApiKey = defaultTmdbApiKey
		return newClient(config, logger)
	}
	return cl, err
}
