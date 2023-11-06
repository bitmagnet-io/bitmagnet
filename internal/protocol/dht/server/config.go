package server

import "time"

type Config struct {
	QueryTimeout time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		QueryTimeout: time.Second * 4,
	}
}
