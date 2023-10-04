package tmdb

import (
	"time"
)

type Config struct {
	ApiKey         string
	RateLimit      time.Duration
	RateLimitBurst int
}

func NewDefaultConfig() Config {
	return Config{
		ApiKey:         defaultTmdbApiKey,
		RateLimit:      defaultRateLimit,
		RateLimitBurst: defaultRateLimitBurst,
	}
}

const (
	defaultTmdbApiKey     = "9c6689fa83ae6814fbfb200d70bba3a8"
	defaultRateLimit      = time.Second / 20
	defaultRateLimitBurst = 5
)
