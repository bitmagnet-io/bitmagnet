package server

import "time"

type Config struct {
	Port         uint16
	QueryTimeout time.Duration
	Addr         string
}

func NewDefaultConfig() Config {
	return Config{
		Port:         3334,
		QueryTimeout: time.Second * 4,
		Addr:         "0.0.0.0",
	}
}
