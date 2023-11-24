package server

import "time"

type Config struct {
	Port         uint16
	QueryTimeout time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		Port:         3334,
		QueryTimeout: time.Second * 4,
	}
}
