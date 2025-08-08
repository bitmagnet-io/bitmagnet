package metrics

import "time"

type Config struct {
	ServiceName       string
	InmemSinkInterval time.Duration
	InmemSinkRetain   time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		ServiceName:       "bitmagnet",
		InmemSinkInterval: time.Minute,
		InmemSinkRetain:   time.Hour * 24 * 30, // 30 days
	}
}
