package tmdb

import "time"

type Config struct {
	Enabled        bool
	BaseUrl        string
	ApiKey         string
	RateLimit      time.Duration
	RateLimitBurst int
}

func NewDefaultConfig() Config {
	return Config{
		Enabled:        true,
		BaseUrl:        "https://api.themoviedb.org/3",
		ApiKey:         defaultTmdbApiKey,
		RateLimit:      defaultRateLimit,
		RateLimitBurst: defaultRateLimitBurst,
	}
}

const (
	defaultTmdbApiKey     = "9c6689fa83ae6814fbfb200d70bba3a8"
	defaultRateLimit      = time.Second / 20
	defaultRateLimitBurst = 8
)
