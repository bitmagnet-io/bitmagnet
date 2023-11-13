package discogs

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
		ApiKey:         defaultDiscogsApiKey,
		RateLimit:      defaultRateLimit,
		RateLimitBurst: defaultRateLimitBurst,
	}
}

const (
	defaultDiscogsApiKey  = ""
	defaultRateLimit      = time.Second / 20
	defaultRateLimitBurst = 5
)
