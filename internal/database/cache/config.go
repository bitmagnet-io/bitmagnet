package cache

import "time"

type Config struct {
	TTL     time.Duration
	MaxKeys uint
}

func NewDefaultConfig() Config {
	return Config{
		TTL:     time.Minute * 10,
		MaxKeys: 1000,
	}
}
