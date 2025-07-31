package tmdb

import "time"

type Config struct {
	BaseURL        string
	APIKey         string
	RateLimit      time.Duration
	RateLimitBurst int
}

func NewDefaultConfig() Config {
	return Config{
		BaseURL:        "https://api.themoviedb.org/3",
		APIKey:         defaultTmdbAPIKey,
		RateLimit:      defaultRateLimit,
		RateLimitBurst: defaultRateLimitBurst,
	}
}

const (
	defaultTmdbAPIKey     = "9c6689fa83ae6814fbfb200d70bba3a8"
	defaultRateLimit      = time.Second / 20
	defaultRateLimitBurst = 5
)
