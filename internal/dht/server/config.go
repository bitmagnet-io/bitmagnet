package server

import "time"

type Config struct {
	QueryTimeout time.Duration
	RateLimit    time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		QueryTimeout: time.Second * 10,
		RateLimit:    time.Second / 100,
	}
}
