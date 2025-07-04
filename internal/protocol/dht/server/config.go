package server

import "time"

type Config struct {
	Port          uint16
	QueryTimeout  time.Duration
	SocketAdapter string
}

func NewDefaultConfig() Config {
	return Config{
		Port:          3334,
		QueryTimeout:  time.Second * 10,
		SocketAdapter: "net",
	}
}
