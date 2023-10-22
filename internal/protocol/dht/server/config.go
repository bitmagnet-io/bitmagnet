package server

import "time"

type Config struct {
	QueryTimeout time.Duration
	RateLimit    time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		QueryTimeout: time.Second * 5,
		RateLimit:    time.Second / 1000,
	}
}
