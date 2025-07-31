package persister

import "time"

type Config struct {
	MaxSize int
	MaxWait time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		MaxSize: 1000,
		MaxWait: 10 * time.Second,
	}
}
