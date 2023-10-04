package metainforequester

import "time"

type Config struct {
	MaxConcurrency uint
	RequestTimeout time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		MaxConcurrency: 100,
		RequestTimeout: 10 * time.Second,
	}
}
