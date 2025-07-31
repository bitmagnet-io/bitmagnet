package metainforequester

import "time"

type Config struct {
	RequestTimeout time.Duration
	MaxConcurrency int
	KeyMutexSize   int
}

func NewDefaultConfig() Config {
	return Config{
		RequestTimeout: 6 * time.Second,
		MaxConcurrency: 1000,
		KeyMutexSize:   1000,
	}
}
