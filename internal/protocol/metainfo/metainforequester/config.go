package metainforequester

import "time"

type Config struct {
	RequestTimeout time.Duration
	KeyMutexSize   int
}

func NewDefaultConfig() Config {
	return Config{
		RequestTimeout: 6 * time.Second,
		KeyMutexSize:   1000,
	}
}
