package server

import "time"

type Config struct {
	Port         uint16
	QueryTimeout time.Duration
	RateLimit    RateLimitConfig
}

type RateLimitConfig struct {
	WindowSize time.Duration
	Burst      int
	Size       int
	TTL        time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		Port:         3334,
		QueryTimeout: time.Second * 4,
		RateLimit: RateLimitConfig{
			WindowSize: time.Second,
			Burst:      4,
			Size:       1000,
			TTL:        time.Second * 20,
		},
	}
}
