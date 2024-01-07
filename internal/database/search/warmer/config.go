package warmer

import "time"

type Config struct {
	Enabled  bool
	Interval time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		Enabled:  true,
		Interval: 50 * time.Minute,
	}
}
