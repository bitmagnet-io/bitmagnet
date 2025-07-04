package cache

import "time"

type Config struct {
	CacheEnabled bool
	TTL          time.Duration
	MaxKeys      uint
}

func NewDefaultConfig() Config {
	return Config{
		CacheEnabled: true,
		TTL:          time.Minute * 10,
		MaxKeys:      1000,
	}
}
