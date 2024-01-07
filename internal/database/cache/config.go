package cache

import "time"

type Config struct {
	CacheEnabled bool
	EaserEnabled bool
	Ttl          time.Duration
	MaxKeys      uint
}

func NewDefaultConfig() Config {
	return Config{
		CacheEnabled: true,
		EaserEnabled: true,
		Ttl:          time.Minute * 60,
		MaxKeys:      1000,
	}
}
