package logging

import (
	"github.com/bitmagnet-io/bitmagnet/internal/logging/filerotator"
	"github.com/bitmagnet-io/bitmagnet/internal/logging/level"
)

type Config struct {
	Level       level.Level
	Development bool
	JSON        bool
	FileRotator filerotator.Config
}

func NewDefaultConfig() Config {
	return Config{
		Level:       "info",
		Development: false,
		JSON:        false,
		FileRotator: filerotator.NewDefaultConfig(),
	}
}
